package topology_test

import (
	"gnet/lib/core"
	"gnet/lib/topology"
	"testing"
	"time"
)

type Game struct {
	*core.ServiceBase
}

func (g *Game) OnRequestMSG(msg *core.Message) {
	g.Respond(msg.Src, core.MSG_ENC_TYPE_GO, msg.Id, "world")
}
func (g *Game) OnCallMSG(msg *core.Message) {
	g.Ret(msg.Src, core.MSG_ENC_TYPE_GO, msg.Id, "world")
}

func (g *Game) OnNormalMSG(msg *core.Message) {
	logsimple.Info("%v", msg)
	//g.RawSend(src, core.MSG_TYPE_NORMAL, "222")
}
func (g *Game) OnDistributeMSG(msg *core.Message) {
	logsimple.Info("%v", msg)
}
func TestMaster(t *testing.T) {
	logsimple.Init("./log", "topology", logsimple.FATAL_LEVEL, logsimple.DEBUG_LEVEL, 10000, 1000)

	core.InitNode(false, true)
	topology.StartMaster("127.0.0.1", "4000")
	core.RegisterNode("./topology")

	game := &Game{core.NewSkeleton(0)}
	id := core.NewService(&core.ModuleParam{
		name: "game1",
		M:    game,
		L:    0,
	})
	logsimple.Info("game1's id :%v", id)

	logsimple.Info("test")
	for {
		time.Sleep(time.Minute * 10)
	}
}
