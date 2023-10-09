package core

import (
	"errors"
	"fmt"
	"gnet/lib/conf"
	"gnet/lib/encoding/gob"
	"gnet/lib/logzap"
	"gnet/lib/timer"
	"gnet/lib/utils"
	"go.uber.org/zap"
	"reflect"
	"sync"
	"time"
)

// 服务接口
type IService interface {
	// 初始化
	OnInit()
	//OnDestory is called when service is closed
	OnDestroy()
	//OnMainLoop is called ever main loop, the delta time is specific by GetDuration()
	OnMainLoop(dt int) //dt is the duration time(unit Millisecond)
	//OnNormalMSG is called when received msg from Send() or RawSend() with MSG_TYPE_NORMAL
	OnNormalMSG(msg *Message)
	//OnSocketMSG is called when received msg from Send() or RawSend() with MSG_TYPE_SOCKET
	OnSocketMSG(msg *Message)
	//OnRequestMSG is called when received msg from Request()
	OnRequestMSG(msg *Message)
	//OnCallMSG is called when received msg from Call()
	OnCallMSG(msg *Message)
	//OnDistributeMSG is called when received msg from Send() or RawSend() with MSG_TYPE_DISTRIBUTE
	OnDistributeMSG(msg *Message)
	//OnCloseNotify is called when received msg from SendClose() with false param.
	OnCloseNotify()
	//~ The very beginning of initializing a service to a new IService
	OnModuleStartup(sid sid, serviceName string)

	getDuration() int
}

// 服务基类, 实现服务接口(IService)
type BaseService struct {
	son               IService                      // 子类/派生类实例
	id                sid                           // 服务ID
	name              string                        // 服务名称
	msgChan           chan *Message                 //
	reqId             uint64                        //
	requestMap        map[uint64]requestCB          //
	requestMutex      sync.Mutex                    //
	callId            uint64                        //
	callChanMap       map[uint64]chan []interface{} //
	callMutex         sync.Mutex                    //
	tick              int                           // 间隔(毫秒ms)
	ts                *timer.TimerSchedule          // 计时器
	normalDispatcher  *CallHelper                   //
	requestDispatcher *CallHelper                   //
	callDispatcher    *CallHelper                   //
}

type requestCB struct {
	respond reflect.Value
	//timeout reflect.Value
}

var (
	ServiceCallTimeout = errors.New("call time out")
)

// 创建服务
func NewService(name string, sz int) *BaseService {
	s := &BaseService{
		son:  nil,
		name: name,
	}
	if sz <= 1024 {
		sz = 1024
	}
	s.msgChan = make(chan *Message, sz)
	s.reqId = 0
	s.requestMap = make(map[uint64]requestCB)
	s.callChanMap = make(map[uint64]chan []interface{})
	return s
}

// 设置子类/派生类
func (s *BaseService) setSon(son IService) {
	s.son = son
}

// 获取服务ID
func (s *BaseService) getId() sid {
	return s.id
}

// 设置服务ID
func (s *BaseService) setId(id sid) {
	s.id = id
}

// 获取服务名称
func (s *BaseService) getName() string {
	return s.name
}

// 设置服务名称
func (s *BaseService) setName(name string) {
	s.name = name
}

func (s *BaseService) pushMSG(m *Message) {
	select {
	case s.msgChan <- m:
	default:
		if s.msgChan == nil {
			logzap.Warn("service error", zap.String("service", s.getName()))
		} else {
			panic(fmt.Sprintf("service is full.<%s>", s.getName()))
		}
	}
}

// [override]销毁
func (s *BaseService) OnDestroy() {
	if s.son != nil {
		s.son.OnDestroy()
		return
	}
	s.destroy()
}

// 销毁
func (s *BaseService) destroy() {
	unregisterService(s)
	msgChan := s.msgChan
	s.msgChan = nil
	close(msgChan)
}

// 启动服务
// @tick  >0时启动计时器
func (s *BaseService) start(tick int) {
	if tick < 0 {
		tick = 0
	}
	s.tick = tick
	if s.tick > 0 && s.ts == nil {
		s.ts = timer.NewTimerSchedule()
		s.ts.SetTick(s.tick)
		s.ts.Start()
	}
	utils.SafeGo(s.loop)
}

// 循环分发消息
func (s *BaseService) loop() {
	// 初始化
	s.OnInit()
	// 循环分发消息
	for {
		if !s.loopSelect() {
			break
		}
	}
	s.OnDestroy()
}

// 分发消息
func (s *BaseService) loopSelect() bool {
	defer func() {
		if err := recover(); err != nil {
			logzap.Error("service error", zap.String("service", s.getName()), zap.String("stack", utils.GetStack()))
		}
	}()
	select {
	case msg, ok := <-s.msgChan:
		if !ok {
			return false
		}
		ok = s.dispatchMSG(msg)
		if !ok {
			return false
		}
	}
	return true
}

// 分发消息
func (s *BaseService) dispatchMSG(msg *Message) bool {
	if msg.EncType == MSG_ENC_TYPE_GO {
		t, err := gob.Unpack(msg.Data[0].([]byte))
		if err != nil {
			panic(err)
		}
		msg.Data = t.([]interface{})
	}
	switch msg.Type {
	case MSG_TYPE_NORMAL:
		s.OnNormalMSG(msg)
	case MSG_TYPE_CLOSE:
		if msg.Data[0].(bool) {
			return false
		}
		s.OnCloseNotify()
	case MSG_TYPE_SOCKET:
		s.OnSocketMSG(msg)
	case MSG_TYPE_REQUEST:
		s.dispatchRequest(msg)
	case MSG_TYPE_RESPOND:
		s.dispatchRespond(msg)
	case MSG_TYPE_CALL:
		s.dispatchCall(msg)
	case MSG_TYPE_DISTRIBUTE:
		s.OnDistributeMSG(msg)
	case MSG_TYPE_TIMEOUT:
		s.dispatchTimeout(msg)
	}
	return true
}

// respndCb is a function like: func(isok bool, ...interface{})  the first param must be a bool
func (s *BaseService) request(dst sid, encType EncType, timeout int, respondCb interface{}, cmd CmdType, data ...interface{}) {
	s.requestMutex.Lock()
	id := s.reqId
	s.reqId++
	cbp := requestCB{reflect.ValueOf(respondCb)}
	s.requestMap[id] = cbp
	s.requestMutex.Unlock()
	utils.PanicWhen(cbp.respond.Kind() != reflect.Func, "respond cb must function.")

	lowLevelSend(s.getId(), dst, MSG_TYPE_REQUEST, encType, id, cmd, data...)

	if timeout > 0 {
		time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
			s.requestMutex.Lock()
			_, ok := s.requestMap[id]
			s.requestMutex.Unlock()
			if ok {
				lowLevelSend(INVALID_SERVICE_ID, s.getId(), MSG_TYPE_TIMEOUT, MSG_ENC_TYPE_NO, id, Cmd_None)
			}
		})
	}
}

func (s *BaseService) dispatchTimeout(m *Message) {
	rid := m.Id
	cbp, ok := s.getDeleteRequestCb(rid)
	if !ok {
		return
	}
	cb := cbp.respond
	var param []reflect.Value
	param = append(param, reflect.ValueOf(true))
	plen := cb.Type().NumIn()
	for i := 1; i < plen; i++ {
		param = append(param, reflect.New(cb.Type().In(i)).Elem())
	}
	cb.Call(param)
}

func (s *BaseService) dispatchRequest(msg *Message) {
	s.OnRequestMSG(msg)
}

func (s *BaseService) respond(dst sid, encType EncType, rid uint64, data ...interface{}) {
	lowLevelSend(s.getId(), dst, MSG_TYPE_RESPOND, encType, rid, Cmd_None, data...)
}

// return request callback by request id
func (s *BaseService) getDeleteRequestCb(id uint64) (requestCB, bool) {
	s.requestMutex.Lock()
	cb, ok := s.requestMap[id]
	delete(s.requestMap, id)
	s.requestMutex.Unlock()
	return cb, ok
}

func (s *BaseService) dispatchRespond(m *Message) {
	var rid uint64
	var data []interface{}
	rid = m.Id
	data = m.Data

	cbp, ok := s.getDeleteRequestCb(rid)
	if !ok {
		return
	}
	cb := cbp.respond
	n := len(data)
	param := make([]reflect.Value, n+1)
	param[0] = reflect.ValueOf(false)
	HelperFunctionToUseReflectCall(cb, param, 1, data)
	cb.Call(param)
}

func (s *BaseService) call(dst sid, encType EncType, cmd CmdType, data ...interface{}) ([]interface{}, error) {
	utils.PanicWhen(dst == s.getId(), "dst must equal to s's id")
	s.callMutex.Lock()
	id := s.callId
	s.callId++
	s.callMutex.Unlock()

	//ch has one buffer, make ret service not block on it.
	ch := make(chan []interface{}, 1)
	s.callMutex.Lock()
	s.callChanMap[id] = ch
	s.callMutex.Unlock()
	if err := lowLevelSend(s.getId(), dst, MSG_TYPE_CALL, encType, id, cmd, data...); err != nil {
		return nil, err
	}
	if conf.CallTimeOut > 0 {
		time.AfterFunc(time.Duration(conf.CallTimeOut)*time.Millisecond, func() {
			s.dispatchRet(id, ServiceCallTimeout)
		})
	}
	ret := <-ch
	s.callMutex.Lock()
	delete(s.callChanMap, id)
	s.callMutex.Unlock()

	close(ch)
	if err, ok := ret[0].(error); ok {
		return ret[1:], err
	}
	return ret, nil
}

func (s *BaseService) callWithTimeout(dst sid, encType EncType, timeout int, cmd CmdType, data ...interface{}) ([]interface{}, error) {
	utils.PanicWhen(dst == s.getId(), "dst must equal to s's id")
	s.callMutex.Lock()
	id := s.callId
	s.callId++
	s.callMutex.Unlock()

	//ch has one buffer, make ret service not block on it.
	ch := make(chan []interface{}, 1)
	s.callMutex.Lock()
	s.callChanMap[id] = ch
	s.callMutex.Unlock()
	if err := lowLevelSend(s.getId(), dst, MSG_TYPE_CALL, encType, id, cmd, data...); err != nil {
		return nil, err
	}
	if timeout > 0 {
		time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
			s.dispatchRet(id, ServiceCallTimeout)
		})
	}
	ret := <-ch
	s.callMutex.Lock()
	delete(s.callChanMap, id)
	s.callMutex.Unlock()

	close(ch)
	if err, ok := ret[0].(error); ok {
		return ret[1:], err
	}
	return ret, nil
}

func (s *BaseService) dispatchCall(msg *Message) {
	s.OnCallMSG(msg)
}

func (s *BaseService) ret(dst sid, encType EncType, cid uint64, data ...interface{}) {
	var dstService *BaseService
	dstService, err := findServiceById(dst)
	if err != nil {
		lowLevelSend(s.getId(), dst, MSG_TYPE_RET, encType, cid, Cmd_None, data...)
		return
	}
	dstService.dispatchRet(cid, data...)
}

func (s *BaseService) dispatchRet(cid uint64, data ...interface{}) {
	s.callMutex.Lock()
	ch, ok := s.callChanMap[cid]
	s.callMutex.Unlock()

	if ok {
		select {
		case ch <- data:
		default:
			utils.PanicWhen(true, "dispatchRet failed on ch.")
		}
	}
}

func (s *BaseService) schedule(interval, repeat int, cb timer.TimerCallback) *timer.Timer {
	utils.PanicWhen(s.tick <= 0, "loopDuraton must greater than zero.")
	return s.ts.Schedule(interval, repeat, cb)
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func (s *BaseService) OnModuleStartup(sid sid, serviceName string) {
	s.normalDispatcher = NewCallHelper(serviceName + ":normalDispatcher")
	s.requestDispatcher = NewCallHelper(serviceName + ":requestDispatcher")
	s.callDispatcher = NewCallHelper(serviceName + ":callDispatcher")
}

// Send调用
func (s *BaseService) Send(dst sid, msgType MsgType, encType EncType, cmd CmdType, data ...interface{}) {
	send(s.getId(), dst, msgType, encType, 0, cmd, data...)
}

// RawSend not encode variables, be careful use
// variables that passed by reference may be changed by others
func (s *BaseService) RawSend(dst sid, msgType MsgType, cmd CmdType, data ...interface{}) {
	sendNoEnc(s.getId(), dst, msgType, 0, cmd, data...)
}

// if isForce is false, then it will just notify the sevice it need to close
// then service can do choose close immediate or close after self clean.
// if isForce is true, then it close immediate
func (s *BaseService) SendClose(dst sid, isForce bool) {
	sendNoEnc(s.getId(), dst, MSG_TYPE_CLOSE, 0, Cmd_None, isForce)
}

// Request send a request msg to dst, and start timeout function if timeout > 0, millisecond
// after receiver call Respond, the responseCb will be called
func (s *BaseService) Request(dst sid, encType EncType, timeout int, responseCb interface{}, cmd CmdType, data ...interface{}) {
	s.request(dst, encType, timeout, responseCb, cmd, data...)
}

// Respond used to respond request msg
func (s *BaseService) Respond(dst sid, encType EncType, rid uint64, data ...interface{}) {
	s.respond(dst, encType, rid, data...)
}

// Call send a call msg to dst, and start a timeout function with the conf.CallTimeOut
// after receiver call Ret, it will return
func (s *BaseService) Call(dst sid, encType EncType, cmd CmdType, data ...interface{}) ([]interface{}, error) {
	return s.call(dst, encType, cmd, data...)
}

// CallWithTimeout send a call msg to dst, and start a timeout function with the timeout millisecond
// after receiver call Ret, it will return
func (s *BaseService) CallWithTimeout(dst sid, encType EncType, timeout int, cmd CmdType, data ...interface{}) ([]interface{}, error) {
	return s.callWithTimeout(dst, encType, timeout, cmd, data...)
}

// Schedule schedule a time with given parameter.
func (s *BaseService) Schedule(interval, repeat int, cb timer.TimerCallback) *timer.Timer {
	if s == nil {
		panic("Schedule must call after OnInit is called(not contain OnInit)")
	}
	return s.schedule(interval, repeat, cb)
}

// Ret used to ret call msg
func (s *BaseService) Ret(dst sid, encType EncType, cid uint64, data ...interface{}) {
	s.ret(dst, encType, cid, data...)
}

func (s *BaseService) OnMainLoop(dt int) {
}
func (s *BaseService) OnNormalMSG(msg *Message) {
	s.normalDispatcher.Call(msg.Cmd, msg.Src, msg.Data...)
}

// [override]初始化
func (s *BaseService) OnInit() {
	if s.son != nil {
		s.son.OnInit()
		return
	}
}

func (s *BaseService) OnSocketMSG(msg *Message) {
}
func (s *BaseService) OnRequestMSG(msg *Message) {
	isAutoReply := s.requestDispatcher.getIsAutoReply(msg.Cmd)
	if isAutoReply { //if auto reply is set, auto respond when user's callback is return.
		ret := s.requestDispatcher.Call(msg.Cmd, msg.Src, msg.Data...)
		s.Respond(msg.Src, msg.EncType, msg.Id, ret...)
	} else { //pass a closure to the user's callback, when to call depends on the user.
		s.requestDispatcher.CallWithReplyFunc(msg.Cmd, msg.Src, func(ret ...interface{}) {
			s.Respond(msg.Src, msg.EncType, msg.Id, ret...)
		}, msg.Data...)
	}
}
func (s *BaseService) OnCallMSG(msg *Message) {
	isAutoReply := s.callDispatcher.getIsAutoReply(msg.Cmd)
	if isAutoReply {
		ret := s.callDispatcher.Call(msg.Cmd, msg.Src, msg.Data...)
		s.Ret(msg.Src, msg.EncType, msg.Id, ret...)
	} else {
		s.callDispatcher.CallWithReplyFunc(msg.Cmd, msg.Src, func(ret ...interface{}) {
			s.Ret(msg.Src, msg.EncType, msg.Id, ret...)
		}, msg.Data...)
	}
}

func (s *BaseService) findCallerByType(msgType MsgType) *CallHelper {
	var caller *CallHelper
	switch msgType {
	case MSG_TYPE_NORMAL:
		caller = s.normalDispatcher
	case MSG_TYPE_REQUEST:
		caller = s.requestDispatcher
	case MSG_TYPE_CALL:
		caller = s.callDispatcher
	default:
		panic("not support msgType")
	}
	return caller
}

// function's first parameter must sid
// isAutoReply: is auto reply when msgType is request or call.
func (s *BaseService) RegisterHandlerFunc(msgType MsgType, cmd CmdType, fun interface{}, isAutoReply bool) {
	caller := s.findCallerByType(msgType)
	caller.AddFunc(cmd, fun)
	caller.setIsAutoReply(cmd, isAutoReply)
}

// method's first parameter must sid
// isAutoReply: is auto reply when msgType is request or call.
func (s *BaseService) RegisterHandlerMethod(msgType MsgType, cmd CmdType, v interface{}, methodName string, isAutoReply bool) {
	caller := s.findCallerByType(msgType)
	caller.AddMethod(cmd, v, methodName)
	caller.setIsAutoReply(cmd, isAutoReply)
}

func (s *BaseService) OnDistributeMSG(msg *Message) {
}

func (s *BaseService) OnCloseNotify() {
	s.SendClose(s.getId(), true)
}
