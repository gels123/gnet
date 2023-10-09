package core_test

import (
	"fmt"
	"gnet/game/conf"
	"gnet/lib/core"
	"gnet/lib/logsimple"
	"testing"
	"time"
)

type Game struct {
	*core.BaseService
	Dst core.sid
}

func (g *Game) OnMainLoop(dt int) {
	g.Send(g.Dst, core.MSG_TYPE_NORMAL, core.MSG_ENC_TYPE_GO, "testNormal", g.Name, []byte{1, 2, 3, 4, 56})
	g.RawSend(g.Dst, core.MSG_TYPE_NORMAL, "testNormal", g.Name, g.Id)

	t := func(timeout bool, data ...interface{}) {
		fmt.Println("request respond ", timeout, data)
	}
	g.Request(g.Dst, core.MSG_ENC_TYPE_GO, 10, t, "testRequest", "hello")

	fmt.Println(g.Call(g.Dst, core.MSG_ENC_TYPE_GO, "testCall", "hello"))
}

func (g *Game) OnInit() {
	//test for go and no enc
	g.RegisterHandlerFunc(core.MSG_TYPE_NORMAL, "testNormal", func(src core.sid, data ...interface{}) {
		logsimple.Info("%v, %v", src, data)
	}, true)
	g.RegisterHandlerFunc(core.MSG_TYPE_REQUEST, "testRequest", func(src core.sid, data ...interface{}) string {
		return "world"
	}, true)
	g.RegisterHandlerFunc(core.MSG_TYPE_CALL, "testCall", func(src core.sid, data ...interface{}) (string, string) {
		return "hello", "world"
	}, true)
}

func TestModule(t *testing.T) {
	logsimple.Init(conf.LogFilePath, conf.LogFileName, conf.LogFileLevel, conf.LogShellLevel, conf.LogMaxLine, conf.LogBufferSize)
	id1 := core.StartService(&core.ModuleParam{
		N: "g1",
		M: &Game{BaseService: core.NewSkeleton(0)},
		L: 0,
	})
	core.StartService(&core.ModuleParam{
		N: "g2",
		M: &Game{BaseService: core.NewSkeleton(1000), Dst: id1},
		L: 0,
	})

	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 1
	}()

	<-ch
}
