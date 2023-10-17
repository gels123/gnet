package core

import (
	"gnet/lib/encoding/gob"
)

type MsgType string
type EncType string
type CmdType string

// 消息结构
type Message struct {
	Src     Sid           // 源服务地址
	Dst     Sid           // 目标服务地址
	Type    MsgType       // 消息类型
	EncType EncType       //
	Id      uint64        // request sid or call sid
	Cmd     CmdType       //
	Data    []interface{} //
}

type NodeInfo struct {
	Name string
	Id   Sid
}

func NewMessage(src, dst Sid, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) *Message {
	if encType == MSG_ENC_TYPE_GOB {
		data = append([]interface{}(nil), gob.Pack(data...))
	}
	msg := &Message{
		Src:     src,
		Dst:     dst,
		Type:    msgType,
		EncType: encType,
		Id:      id,
		Cmd:     cmd,
		Data:    data,
	}
	return msg
}

func init() {
	gob.RegisterStructType(Message{})
	gob.RegisterStructType(NodeInfo{})
}

func sendNoEnc(src Sid, dst Sid, msgType MsgType, id uint64, cmd CmdType, data ...interface{}) error {
	return lowLevelSend(src, dst, msgType, MSG_ENC_TYPE_NO, id, cmd, data...)
}

func send(src Sid, dst Sid, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) error {
	return lowLevelSend(src, dst, msgType, encType, id, cmd, data...)
}

func lowLevelSend(src, dst Sid, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) error {
	dsts, err := findServiceById(dst)
	isLocal := isLocalSid(dst)
	//a local service is not been found
	if err != nil && isLocal {
		return err
	}
	var msg *Message
	msg = NewMessage(src, dst, msgType, encType, id, cmd, data...)
	if err != nil {
		//doesn't find service and dstid is remote sid, send a forward msg to router.
		route(Cmd_Forward, msg)
		return nil
	}
	dsts.pushMsg(msg)
	return nil
}

// send msg to dst by dst's ServiceBase name
func sendName(src Sid, dst string, msgType MsgType, cmd CmdType, data ...interface{}) error {
	dsts, err := findServiceByName(dst)
	if err != nil {
		return err
	}
	return lowLevelSend(src, dsts.getId(), msgType, MSG_ENC_TYPE_GOB, 0, cmd, data...)
}

// ForwardLocal forward the message to the specified local sevice.
func ForwardLocal(m *Message) {
	dsts, err := findServiceById(Sid(m.Dst))
	if err != nil {
		return
	}
	switch m.Type {
	case MSG_TYPE_NORMAL,
		MSG_TYPE_REQUEST,
		MSG_TYPE_RESPOND,
		MSG_TYPE_CALL,
		MSG_TYPE_DISTRIBUTE:
		dsts.pushMsg(m)
	case MSG_TYPE_RET:
		if m.EncType == MSG_ENC_TYPE_GOB {
			t, err := gob.Unpack(m.Data[0].([]byte))
			if err != nil {
				panic(err)
			}
			m.Data = t.([]interface{})
		}
		cid := m.Id
		dsts.dispatchRet(cid, m.Data...)
	}
}

// DistributeMSG distribute the message to all local sevice
func DistributeMSG(src Sid, cmd CmdType, data ...interface{}) {
	mgr.dicMutex.Lock()
	defer mgr.dicMutex.Unlock()
	for dst, ser := range mgr.dictId {
		if Sid(dst) != src {
			localSendWithoutMutex(src, ser, MSG_TYPE_DISTRIBUTE, MSG_ENC_TYPE_NO, 0, cmd, data...)
		}
	}
}

// localSendWithoutMutex send a message to the local ServiceBase with no mutex.
func localSendWithoutMutex(src Sid, dstService *ServiceBase, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) {
	msg := NewMessage(src, dstService.getId(), msgType, encType, id, cmd, data...)
	dstService.pushMsg(msg)
}
