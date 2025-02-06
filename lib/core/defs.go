/*
 * 服务相关定义
 */
package core

// 服务ID规则配置（高16位为集群节点ID, 低48位为服务ID）
const (
	NODE_ID_OFF          = 64 - 16
	NODE_ID_MAX          = 0xFFFF
	NODE_ID_MASK         = 0xFFFF << NODE_ID_OFF
	INVALID_SRC_ID       = NODE_ID_MASK
	MASTER_NODE_ID       = 0
	MIN_SRC_ID SID 	 	 = 10
)

// 服务ID定义(高16位为节点ID,低48位为服务ID)
type SID uint64

// 集群节点ID
func (sid SID) NodeId() uint64 {
	return (uint64(sid) & NODE_ID_MASK) >> NODE_ID_OFF
}

// 服务ID
func (sid SID) BaseId() uint64 {
	return uint64(sid) & (^uint64(NODE_ID_MASK))
}

// 是否合法
func (sid SID) Valid() bool {
	return !(sid == INVALID_SRC_ID || sid == 0)
}

// 是否不合法
func (sid SID) Invalid() bool {
	return sid == INVALID_SRC_ID || sid == 0
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
消息类型
*/
const (
	MSG_TYPE_NORMAL     MsgType = "Msg.Normal"
	MSG_TYPE_REQUEST    MsgType = "Msg.Request"
	MSG_TYPE_RESPOND    MsgType = "Msg.Respond"
	MSG_TYPE_TIMEOUT    MsgType = "Msg.TimeOut"
	MSG_TYPE_CALL       MsgType = "Msg.Call"
	MSG_TYPE_RET        MsgType = "Msg.Ret"
	MSG_TYPE_CLOSE      MsgType = "Msg.Close"
	MSG_TYPE_SOCKET     MsgType = "Msg.Socket"
	MSG_TYPE_ERR        MsgType = "Msg.Error"
	MSG_TYPE_DISTRIBUTE MsgType = "Msg.Distribute"
)

/*
消息编码类型
*/
const (
	MSG_ENC_TYPE_NIL EncType = "nil"
	MSG_ENC_TYPE_GOB EncType = "gob"
)
