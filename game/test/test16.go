package main

import (
	"gnet/lib/core"
	"gnet/lib/logsimple"
	"time"
)

type Game struct {
	*core.ServiceBase
	Dst core.sid
}

//func (g *Game) OnMainLoop(dt int) {
//	g.Send(g.Dst, core.MSG_TYPE_NORMAL, core.MSG_ENC_TYPE_GOB, "testNormal", g.Name, []byte{1, 2, 3, 4, 56})
//	g.RawSend(g.Dst, core.MSG_TYPE_NORMAL, "testNormal", g.Name, g.Id)
//
//	t := func(timeout bool, data ...interface{}) {
//		fmt.Println("request respond ", timeout, data)
//	}
//	g.Request(g.Dst, core.MSG_ENC_TYPE_GOB, 10, t, "testRequest", "hello")
//
//	fmt.Println(g.Call(g.Dst, core.MSG_ENC_TYPE_GOB, "testCall", "hello"))
//}
//
//func (g *Game) OnInit() {
//	//test for go and no enc
//	g.RegisterHandlerFunc(core.MSG_TYPE_NORMAL, "testNormal", func(src core.ServiceID, data ...interface{}) {
//		log.Info("%v, %v", src, data)
//	}, true)
//	g.RegisterHandlerFunc(core.MSG_TYPE_REQUEST, "testRequest", func(src core.ServiceID, data ...interface{}) string {
//		return "world"
//	}, true)
//	g.RegisterHandlerFunc(core.MSG_TYPE_CALL, "testCall", func(src core.ServiceID, data ...interface{}) (string, string) {
//		return "hello", "world"
//	}, true)
//}

func main() {
	//logsimple.Init(conf.LogFilePath, conf.LogFileName, conf.LogFileLevel, conf.LogShellLevel, conf.LogMaxLine, conf.LogBufferSize)
	logsimple.Info("===== test16 main =====")
	//id1 := core.NewService(&core.ModuleParam{
	//	name: "g1",
	//	M: &Game{Skeleton: core.NewSkeleton(0)},
	//	L: 0,
	//})
	//core.NewService(&core.ModuleParam{
	//	name: "g2",
	//	M: &Game{Skeleton: core.NewSkeleton(1000), Dst: id1},
	//	L: 0,
	//})

	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 1
	}()

	<-ch
}
