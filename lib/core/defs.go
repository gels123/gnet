package core

// 服务ID规则配置（高16位为集群节点ID, 低48位为服务ID）
const (
	NODE_ID_OFF          = 64 - 16
	NODE_ID_MAX          = 0xFFFF
	NODE_ID_MASK         = 0xFFFF << NODE_ID_OFF
	INVALID_SRC_ID       = NODE_ID_MASK
	MASTER_NODE_ID       = 0
	MIN_SRC_ID     SvcId = 10
)

// 服务ID定义(高16位为节点ID,低48位为服务ID)
type SvcId uint64

// 集群节点ID
func (id SvcId) NodeId() uint64 {
	return (uint64(id) & NODE_ID_MASK) >> NODE_ID_OFF
}

// 服务ID
func (id SvcId) BaseId() uint64 {
	return uint64(id) & (^uint64(NODE_ID_MASK))
}

// 是否合法
func (id SvcId) Valid() bool {
	return !(id == INVALID_SRC_ID || id == 0)
}

// 是否不合法
func (id SvcId) Invalid() bool {
	return id == INVALID_SRC_ID || id == 0
}

const (
	Cmd_None                    CmdType = "CmdType.Core.None"
	Cmd_Forward                 CmdType = "CmdType.Core.Forward"
	Cmd_Distribute              CmdType = "CmdType.Core.Distribute"
	Cmd_RegisterNode            CmdType = "CmdType.Core.RegisterNode"
	Cmd_RegisterNodeRet         CmdType = "CmdType.Core.RegisterNodeRet"
	Cmd_RegisterName            CmdType = "CmdType.Core.RegisterName"
	Cmd_GetIdByName             CmdType = "CmdType.Core.GetIdByName"
	Cmd_GetIdByNameRet          CmdType = "CmdType.Core.GetIdByNameRet"
	Cmd_NameAdd                 CmdType = "CmdType.Core.NameAdd"
	Cmd_NameDeleted             CmdType = "CmdType.Core.NameDeleted"
	Cmd_Exit                    CmdType = "CmdType.Core.Exit"
	Cmd_Exit_Node               CmdType = "CmdType.Core.ExitNode"
	Cmd_Default                 CmdType = "CmdType.Core.Default"
	Cmd_RefreshSlaveWhiteIPList CmdType = "CmdType.Core.RefreshSlaveWhiteIPList"
)

/*
const (

	MSG_TYPE_NORMAL = iota
	MSG_TYPE_REQUEST
	MSG_TYPE_RESPOND
	MSG_TYPE_TIMEOUT
	MSG_TYPE_CALL
	MSG_TYPE_RET
	MSG_TYPE_CLOSE
	MSG_TYPE_SOCKET
	MSG_TYPE_ERR
	MSG_TYPE_DISTRIBUTE

)
*/
const (
	MSG_TYPE_NORMAL     MsgType = "MsgType.Normal"
	MSG_TYPE_REQUEST    MsgType = "MsgType.Request"
	MSG_TYPE_RESPOND    MsgType = "MsgType.Respond"
	MSG_TYPE_TIMEOUT    MsgType = "MsgType.TimeOut"
	MSG_TYPE_CALL       MsgType = "MsgType.Call"
	MSG_TYPE_RET        MsgType = "MsgType.Ret"
	MSG_TYPE_CLOSE      MsgType = "MsgType.Close"
	MSG_TYPE_SOCKET     MsgType = "MsgType.Socket"
	MSG_TYPE_ERR        MsgType = "MsgType.Error"
	MSG_TYPE_DISTRIBUTE MsgType = "MsgType.Distribute"
)

/*
消息编码类型定义
const (

	MSG_ENC_TYPE_NIL = iota
	MSG_ENC_TYPE_GOB

)
*/
const (
	MSG_ENC_TYPE_NIL EncType = "nil"
	MSG_ENC_TYPE_GOB EncType = "gob"
)
