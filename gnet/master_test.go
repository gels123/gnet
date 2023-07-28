package lotou_test

import (
	"github.com/sydnash/lotou"
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"testing"
)

type Game struct {
	*core.Skeleton
	remoteId core.ServiceID
}

/*
func (g *Game) OnRequestMSG(src core.ServiceID, rid uint64, data ...interface{}) {
	g.Respond(src, rid, "world")
}
func (g *Game) OnCallMSG(src core.ServiceID, cid uint64, data ...interface{}) {
	g.Ret(src, cid, "world")
}

func (g *Game) OnNormalMSG(src core.ServiceID, data ...interface{}) {
	log.Info("%v, %v", src, data)
	//g.RawSend(src, core.MSG_TYPE_NORMAL, "222")
}*/
func (g *Game) OnDistributeMSG(msg *core.Message) {
	log.Info("%v", msg)
}
func (g *Game) OnInit() {
	log.Info("OnInit: name:%v  id:%v", g.Name, g.Id)
	g.remoteId, _ = core.NameToId("game1")
	log.Info("name2id: game1:%v", g.remoteId)
	if g.D > 0 {
		g.Schedule(100, 0, func(dt int) {
			log.Info("time schedule.")
		})
	}
	g.RegisterHandlerFunc(core.MSG_TYPE_NORMAL, "testNormal", func(src core.ServiceID, data ...interface{}) {
		log.Info("%v, %v", src, data)
	}, true)
	g.RegisterHandlerFunc(core.MSG_TYPE_REQUEST, "testRequest", func(src core.ServiceID, data ...interface{}) string {
		return "world"
	}, true)
	g.RegisterHandlerFunc(core.MSG_TYPE_CALL, "testCall", func(src core.ServiceID, data ...interface{}) (string, string) {
		return "hello", "world"
	}, true)
}

func TestMaster(t *testing.T) {
	conf.CoreIsStandalone = false
	conf.CoreIsMaster = true
	game := &Game{core.NewSkeleton(0), 0}
	lotou.Start(nil, &core.ModuleParam{
		N: "game1",
		M: game,
		L: 0,
	})
}
