package core

import (
	"errors"
	"fmt"
	"gnet/lib/logzap"
	"gnet/lib/utils"
	"reflect"
)

// CallHelper use reflect.Call to invoke a function.
// it's not thread safe
const (
	ReplyFuncPosition = 1
)

// 错误定义
var (
	ErrFuncNotFound = errors.New("func not found.")
)

type callbackDesc struct {
	cb          reflect.Value
	isAutoReply bool
}

type CallHelper struct {
	serviceName string //help to locate which callback is not registered.
	funcMap     map[CmdType]*callbackDesc
}

type ReplyFunc func(data ...interface{})

func NewCallHelper(name string) *CallHelper {
	return &CallHelper{
		serviceName: name,
		funcMap:     make(map[CmdType]*callbackDesc),
	}
}

// AddFunc add callback with normal function
func (c *CallHelper) AddFunc(cmd CmdType, fun interface{}) {
	f := reflect.ValueOf(fun)
	utils.PanicWhen(f.Kind() != reflect.Func, "fun must be a function type.")
	c.funcMap[cmd] = &callbackDesc{f, true}
}

// AddMethod add callback with struct's method by method name. The method name must be exported.
func (c *CallHelper) AddMethod(cmd CmdType, v interface{}, methodName string) {
	val := reflect.ValueOf(v)
	f := val.MethodByName(methodName)
	utils.PanicWhen(f.Kind() != reflect.Func, fmt.Sprintf("[CallHelper:AddMethod] cmd{%v} method must be a function type.", cmd))
	c.funcMap[cmd] = &callbackDesc{f, true}
}

// setIsAutoReply recode special cmd is auto reply after Call is return
func (c *CallHelper) setIsAutoReply(cmd CmdType, isAutoReply bool) {
	cb := c.findCallbackDesc(cmd)
	cb.isAutoReply = isAutoReply
	if !isAutoReply {
		t := reflect.New(cb.cb.Type().In(ReplyFuncPosition))
		_ = t.Elem().Interface().(ReplyFunc)
	}
}

func (c *CallHelper) getIsAutoReply(cmd CmdType) bool {
	return c.findCallbackDesc(cmd).isAutoReply
}

// 查找函数
func (c *CallHelper) findCallbackDesc(cmd CmdType) *callbackDesc {
	cb, ok := c.funcMap[cmd]
	if !ok {
		if cb, ok = c.funcMap[Cmd_Default]; ok {

		} else {
			logzap.Panicw("func not found", "serviceName", c.serviceName, "cmd", cmd)
		}
	}
	return cb
}

// Call invoke special function for cmd
func (c *CallHelper) Call(cmd CmdType, src SvcId, param ...interface{}) []interface{} {
	defer func() {
		if err := recover(); err != nil {
			logzap.Errorw("CallHelper.Call error", "serviceName", c.serviceName, "cmd", cmd, "err", err)
		}
	}()

	cb := c.findCallbackDesc(cmd)
	//addition one param for source service sid
	p := make([]reflect.Value, len(param)+1)
	p[0] = reflect.ValueOf(src) //append src service sid
	HelperFunctionToUseReflectCall(cb.cb, p, 1, param)

	ret := cb.cb.Call(p)

	out := make([]interface{}, len(ret))
	for i, v := range ret {
		out[i] = v.Interface()
	}
	return out
}

// CallWithReplyFunc invoke special function for cmd with a reply function which is used to reply Call or Request.
func (c *CallHelper) CallWithReplyFunc(cmd CmdType, src SvcId, replyFunc ReplyFunc, param ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logzap.Errorw("CallHelper.Call error", "serviceName", c.serviceName, "cmd", cmd, "err", err)
		}
	}()

	cb := c.findCallbackDesc(cmd)
	//addition two param for source service sid and reply function
	p := make([]reflect.Value, len(param)+2)
	p[0] = reflect.ValueOf(src) //append src service sid
	p[1] = reflect.ValueOf(replyFunc)
	HelperFunctionToUseReflectCall(cb.cb, p, 2, param)

	cb.cb.Call(p)
}
