/*
 * 消息
 */
package core

import (
	"gnet/lib/encoding/gob"
)

type MsgType string
type EncType string
type CmdType string

// 消息结构
type Message struct {
	Src     SID         // 源地址
	Dst     SID         // 目标地址
	Type    MsgType       // 消息类型
	EncType EncType       // 序列化类型
	Id      uint64        // request session or call session id
	Cmd     CmdType       // 函数指令
	Data    []interface{} // 函数参数
}

type NodeInfo struct {
	Name string
	Id   SID
}

func NewMessage(src, dst SID, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) *Message {
	if encType == MSG_ENC_TYPE_GOB {
		data = append([]interface{}(nil), gob.Pack(data...))
	}
	msg := &Message {
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

func send(src, dst SID, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) error {
	isLocal := isLocalSid(dst)
	service, err := findServiceById(dst)
	//Local service not found
	if isLocal && err != nil {
		return err
	}
	var msg *Message
	msg = NewMessage(src, dst, msgType, encType, id, cmd, data...)
	if err != nil {
		//doesn't find service and dstid is remote sid, send a forward msg to router.
		route(Cmd_Forward, msg)
		return nil
	}
	service.pushMsg(msg)
	return nil
}

func sendRaw(src SID, dst SID, msgType MsgType, id uint64, cmd CmdType, data ...interface{}) error {
	return send(src, dst, msgType, MSG_ENC_TYPE_NIL, id, cmd, data...)
}

func sendByName(src SID, dst string, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) error {
	service, err := findServiceByName(dst)
	if err != nil {
		return err
	}
	return send(src, service.GetId(), msgType, encType, id, cmd, data...)
}

// ForwardLocal forward the message to the specified local sevice.
func ForwardLocal(msg *Message) {
	dsts, err := findServiceById(SID(msg.Dst))
	if err != nil {
		return
	}
	switch msg.Type {
	case MSG_TYPE_NORMAL,
		MSG_TYPE_REQUEST,
		MSG_TYPE_RESPOND,
		MSG_TYPE_CALL,
		MSG_TYPE_DISTRIBUTE:
		dsts.pushMsg(msg)
	case MSG_TYPE_RET:
		if msg.EncType == MSG_ENC_TYPE_GOB {
			t, err := gob.Unpack(msg.Data[0].([]byte))
			if err != nil {
				panic(err)
			}
			msg.Data = t.([]interface{})
		}
		cid := msg.Id
		dsts.dispatchRet(cid, msg.Data...)
	}
}

// DistributeMSG distribute the message to all local sevice
func DistributeMSG(src SID, cmd CmdType, data ...interface{}) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	for dst, s := range mgr.sidDict {
		if SID(dst) != src {
			sendLocal(src, s, MSG_TYPE_DISTRIBUTE, MSG_ENC_TYPE_NIL, 0, cmd, data...)
		}
	}
}

// sendLocal send a message to the local ServiceBase with no mutex.
func sendLocal(src SID, dst *ServiceBase, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) {
	var msg *Message
	msg = NewMessage(src, dst.GetId(), msgType, encType, id, cmd, data...)
	dst.pushMsg(msg)
}

func init() {
	gob.RegisterStructType(Message{})
	gob.RegisterStructType(NodeInfo{})
}
