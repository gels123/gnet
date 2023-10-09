package core

import (
	"gnet/lib/encoding/gob"
)

type MsgType string
type EncType string
type CmdType string

// Message is the based struct of msg through all service
// by convention, the first value of Data is a string as the method name
type Message struct {
	Src     sid
	Dst     sid
	Type    MsgType // Used to be int32
	EncType EncType
	Id      uint64 //request id or call id
	Cmd     CmdType
	Data    []interface{}
}

type NodeInfo struct {
	Name string
	Id   sid
}

func NewMessage(src, dst sid, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) *Message {
	switch encType {
	case MSG_ENC_TYPE_NO:
	case MSG_ENC_TYPE_GO:
		data = append([]interface{}(nil), gob.Pack(data...))
	}
	msg := &Message{src, dst, msgType, encType, id, cmd, data}
	return msg
}

func init() {
	gob.RegisterStructType(Message{})
	gob.RegisterStructType(NodeInfo{})
}

func sendNoEnc(src sid, dst sid, msgType MsgType, id uint64, cmd CmdType, data ...interface{}) error {
	return lowLevelSend(src, dst, msgType, MSG_ENC_TYPE_NO, id, cmd, data...)
}

func send(src sid, dst sid, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) error {
	return lowLevelSend(src, dst, msgType, encType, id, cmd, data...)
}

func lowLevelSend(src, dst sid, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) error {
	dsts, err := findServiceById(dst)
	isLocal := checkIsLocalId(dst)
	//a local service is not been found
	if err != nil && isLocal {
		return err
	}
	var msg *Message
	msg = NewMessage(src, dst, msgType, encType, id, cmd, data...)
	if err != nil {
		//doesn't find service and dstid is remote id, send a forward msg to router.
		route(Cmd_Forward, msg)
		return nil
	}
	dsts.pushMSG(msg)
	return nil
}

// send msg to dst by dst's BaseService name
func sendName(src sid, dst string, msgType MsgType, cmd CmdType, data ...interface{}) error {
	dsts, err := findServiceByName(dst)
	if err != nil {
		return err
	}
	return lowLevelSend(src, dsts.getId(), msgType, MSG_ENC_TYPE_GO, 0, cmd, data...)
}

// ForwardLocal forward the message to the specified local sevice.
func ForwardLocal(m *Message) {
	dsts, err := findServiceById(sid(m.Dst))
	if err != nil {
		return
	}
	switch m.Type {
	case MSG_TYPE_NORMAL,
		MSG_TYPE_REQUEST,
		MSG_TYPE_RESPOND,
		MSG_TYPE_CALL,
		MSG_TYPE_DISTRIBUTE:
		dsts.pushMSG(m)
	case MSG_TYPE_RET:
		if m.EncType == MSG_ENC_TYPE_GO {
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
func DistributeMSG(src sid, cmd CmdType, data ...interface{}) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	for dst, ser := range h.dic {
		if sid(dst) != src {
			localSendWithoutMutex(src, ser, MSG_TYPE_DISTRIBUTE, MSG_ENC_TYPE_NO, 0, cmd, data...)
		}
	}
}

// localSendWithoutMutex send a message to the local BaseService with no mutex.
func localSendWithoutMutex(src sid, dstService *BaseService, msgType MsgType, encType EncType, id uint64, cmd CmdType, data ...interface{}) {
	msg := NewMessage(src, dstService.getId(), msgType, encType, id, cmd, data...)
	dstService.pushMSG(msg)
}
