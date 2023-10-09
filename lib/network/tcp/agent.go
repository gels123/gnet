package tcp

import (
	"bufio"
	"gnet/lib/core"
	"gnet/lib/utils"
	"net"
	"time"
)

/*
	recieve message from network
	split message into package
	send package to dest service

	reciev other inner message and process(such as write to network; close; change dest service)

	agent has two goroutine:
		one to read message from tcpConnector and send to dest
		one read inner message chan and process

	this also has a timeout for first data arrived after agent create.
*/

type Agent struct {
	*core.BaseService
	Con                  *net.TCPConn
	closeing             bool
	hostService          core.sid
	hasDataArrived       bool
	leftTimeBeforArrived int
	inbuffer             []byte
	outbuffer            *bufio.Writer
	parseCache           *ParseCache
}

const AgentNoDataHoldtime = 5000

func (a *Agent) OnInit() {
	a.hasDataArrived = false
	a.leftTimeBeforArrived = AgentNoDataHoldtime
	a.inbuffer = make([]byte, DEFAULT_RECV_BUFF_LEN)
	a.outbuffer = bufio.NewWriter(a.Con)
	a.parseCache = &ParseCache{}
	logsimple.Info("receive a connect at: %v->%v", a.Con.LocalAddr(), a.Con.RemoteAddr())
	a.sendToHost(core.MSG_TYPE_SOCKET, AGENT_ARRIVE, a.Con.LocalAddr().String(), a.Con.RemoteAddr().String())
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logsimple.Error("recover: stack: %v\n, %v", utils.GetStack(), err)
			}
		}()
		for {
			pack, err := Subpackage(a.inbuffer, a.Con, a.parseCache)
			if err != nil {
				logsimple.Error("agent read msg failed: %s", err)
				a.onConnectError()
				break
			}
			if !a.hasDataArrived {
				a.hasDataArrived = true
			}
			for _, v := range pack {
				a.sendToHost(core.MSG_TYPE_SOCKET, AGENT_DATA, v)
			}
		}
	}()
}

func NewAgent(con *net.TCPConn, hostID core.sid) *Agent {
	a := &Agent{
		Con:         con,
		hostService: hostID,
		BaseService: core.NewSkeleton(5000),
	}
	return a
}

func (a *Agent) OnMainLoop(dt int) {
	if !a.hasDataArrived {
		a.leftTimeBeforArrived -= dt
		if a.leftTimeBeforArrived < 0 {
			logsimple.Error("agent hasn't got any data within %v ms", AgentNoDataHoldtime)
			a.SendClose(a.Id, false)
		}
	}
}

func (a *Agent) OnNormalMSG(msg *core.Message) {
	if a.closeing {
		return
	}
	data := msg.Data
	cmd := msg.Cmd
	if cmd == AGENT_CMD_SEND {
		a.Con.SetWriteDeadline(time.Now().Add(time.Second * 20))
		msg := data[0].([]byte)
		if _, err := a.outbuffer.Write(msg); err != nil {
			logsimple.Error("agent write msg failed: %s", err)
			a.onConnectError()
			return
		}
		if err := a.outbuffer.Flush(); err != nil {
			logsimple.Error("agent flush msg failed: %s", err)
			a.onConnectError()
			return
		}
	}
}

func (a *Agent) OnDestroy() {
	a.close()
}

func (a *Agent) onConnectError() {
	a.sendToHost(core.MSG_TYPE_SOCKET, AGENT_CLOSED)
	a.SendClose(a.Id, false)
	a.closeing = true
}

func (a *Agent) close() {
	logsimple.Info("close agent. %v", a.Con.RemoteAddr())
	a.Con.Close()
}

func (a *Agent) sendToHost(msgType core.MsgType, cmd core.CmdType, data ...interface{}) {
	a.RawSend(a.hostService, msgType, cmd, data...)
}
