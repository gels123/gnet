package main

import (
	"fmt"
	"gnet/lib/core"
	"time"
)

// type Game struct {
// 	*core.ServiceBase
// 	Dst core.SID
// }

// func (g *Game) OnMainLoop(dt int) {
// 	g.Send(g.Dst, core.MSG_TYPE_NORMAL, core.MSG_ENC_TYPE_GO, "testNormal", g.Name, []byte{1, 2, 3, 4, 56})
// 	g.RawSend(g.Dst, core.MSG_TYPE_NORMAL, "testNormal", g.Name, g.Id)

// 	t := func(timeout bool, data ...interface{}) {
// 		fmt.Println("request respond ", timeout, data)
// 	}
// 	g.Request(g.Dst, core.MSG_ENC_TYPE_GO, 10, t, "testRequest", "hello")

// 	fmt.Println(g.Call(g.Dst, core.MSG_ENC_TYPE_GO, "testCall", "hello"))
// }

// func (g *Game) OnInit() {
// 	//test for go and no enc
// 	g.RegisterHandlerFunc(core.MSG_TYPE_NORMAL, "testNormal", func(src core.Sid, data ...interface{}) {
// 		logsimple.Info("%v, %v", src, data)
// 	}, true)
// 	g.RegisterHandlerFunc(core.MSG_TYPE_REQUEST, "testRequest", func(src core.Sid, data ...interface{}) string {
// 		return "world"
// 	}, true)
// 	g.RegisterHandlerFunc(core.MSG_TYPE_CALL, "testCall", func(src core.Sid, data ...interface{}) (string, string) {
// 		return "hello", "world"
// 	}, true)
// }

func main() {
	fmt.Println("==main begin==")

	opt := core.ServiceOption {
		Name:  "test",
		MsgSz: 1024,
		Tick:  0,
	}
	fmt.Println("create service opt", opt)

	// svc1 := core.NewService(&core.ServiceOption {
	// 	name: "svc1",
	// 	msgSz: 1024,
	// 	tick: 0,
	// })
	// core.NewService(&core.ModuleParam{
	// 	name: "g2",
	// 	M:    &Game{ServiceBase: core.NewSkeleton(1000), Dst: svc1},
	// 	L:    0,
	// })

	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 1
	}()
	<-ch

	fmt.Println("==main end==")
}