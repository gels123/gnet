package core

import (
	"fmt"
	"gnet/lib/logsimple"
	"testing"
	"time"
)

// 测试服务
type TestService struct {
	base *core.ServiceBase //
}

func newTestService() *Doctor {
	d := &TestService{
		base: core.NewService("testsvc", 1024),
	}
	d.base.setSelf(d)
	return d
}

func (g *TestService) OnMainLoop(dt int) {
	g.Send(g.Dst, core.MSG_TYPE_NORMAL, core.MSG_ENC_TYPE_GO, "testNormal", g.Name, []byte{1, 2, 3, 4, 56})
	g.RawSend(g.Dst, core.MSG_TYPE_NORMAL, "testNormal", g.Name, g.Id)

	t := func(timeout bool, data ...interface{}) {
		fmt.Println("request respond ", timeout, data)
	}
	g.Request(g.Dst, core.MSG_ENC_TYPE_GO, 10, t, "testRequest", "hello")

	fmt.Println(g.Call(g.Dst, core.MSG_ENC_TYPE_GO, "testCall", "hello"))
}

func (g *TestService) OnInit() {
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
	testService := core.NewService(&core.ModuleParam{
		name: "g1",
		M:    &TestService{ServiceBase: core.NewSkeleton(0)},
		L:    0,
	})

	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 1
	}()

	<-ch
}
