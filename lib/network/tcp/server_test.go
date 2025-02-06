package tcp_test

import (
	"gnet/lib/core"
	"gnet/lib/encoding/binary"
	"gnet/lib/network/tcp"
	"testing"
)

type M struct {
	*core.ServiceBase
	decoder *binary.Decoder
}

func (m *M) OnNormalMSG(msg *core.Message) {
	cmd := msg.Cmd
	if cmd == tcp.AGENT_CLOSED {
		logsimple.Info("agent closed")
	}
}

func (m *M) OnSocketMSG(msg *core.Message) {
	src := msg.From
	cmd := msg.Cmd
	data := msg.Data
	if cmd == tcp.AGENT_DATA {
		data := data[0].([]byte)
		m.decoder.SetBuffer(data)
		var msg []byte = []byte{}
		m.decoder.Decode(&msg)
		logsimple.Info("%v, %v", src, string(msg))

		m.RawSend(src, core.MSG_TYPE_NORMAL, tcp.AGENT_CMD_SEND, data)
	}
}

func TestServer(t *testing.T) {
	logsimple.Init("./log", "tcpserver", logsimple.FATAL_LEVEL, logsimple.DEBUG_LEVEL, 10000, 1000)
	m := &M{ServiceBase: core.NewSkeleton(0)}
	m.decoder = binary.NewDecoder()
	core.NewService(&core.ModuleParam{
		name: ".m",
		M:    m,
		L:    0,
	})

	s := tcp.NewServer("", "3333", m.Id)
	s.Listen()

	ch := make(chan int)
	<-ch

	s.Close()
}
