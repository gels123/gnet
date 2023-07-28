package topology_test

import (
	"gnet/lib/core"
	"gnet/lib/log"
	"gnet/lib/topology"
	"testing"
	"time"
)

type Game struct {
	*core.Skeleton
}

func (g *Game) OnRequestMSG(msg *core.Message) {
	g.Respond(msg.Src, core.MSG_ENC_TYPE_GO, msg.Id, "world")
}
func (g *Game) OnCallMSG(msg *core.Message) {
	g.Ret(msg.Src, core.MSG_ENC_TYPE_GO, msg.Id, "world")
}

func (g *Game) OnNormalMSG(msg *core.Message) {
	log.Info("%v", msg)
	//g.RawSend(src, core.MSG_TYPE_NORMAL, "222")
}
func (g *Game) OnDistributeMSG(msg *core.Message) {
	log.Info("%v", msg)
}
func TestMaster(t *testing.T) {
	log.Init("./log", "topology", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)

	core.InitNode(false, true)
	topology.StartMaster("127.0.0.1", "4000")
	core.RegisterNode("./topology")

	game := &Game{core.NewSkeleton(0)}
	id := core.StartService(&core.ModuleParam{
		N: "game1",
		M: game,
		L: 0,
	})
	log.Info("game1's id :%v", id)

	log.Info("test")
	for {
		time.Sleep(time.Minute * 10)
	}
}
