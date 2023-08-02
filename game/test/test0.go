package main

import (
	"fmt"
	"time"
)

type Student_ struct {
	name string
	age  int
}

func logf1() {
	//log.Init("../../log", "game", log.DEBUG_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	//s := &Student_{"yyyyy", 100}
	//log.Debug("hahaha %v, %v", 2, s)
	//log.Error("hahaha %v, %v", 2, s)
	//log.Warn("hahaha %v, %v", 2, s)
	////log.Fatal("hahaha %v, %v", 2, s)
}

func log_zap() {
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	//log.Print("hello world")

	//logger, _ := zap.NewProduction()
	//defer logger.Sync() // flushes buffer, if any
	//sugar := logger.Sugar()
	//sugar.Infow("failed to fetch URL",
	//	// Structured context as loosely typed key-value pairs.
	//	"url", "wwwwwffff",
	//	"attempt", 3,
	//	"backoff", time.Second,
	//)
	//sugar.Infof("Failed to fetch URL: %s", "wwwwwffff")
	//logger.Info("my test log", zap.Bool("player", true), zap.String("name", "llily"))
	//
	//logzap.Infof("Failed to fetch URL: %s", "xxxx1")
	//logzap.Errorf("Failed to fetch URL: %s", "xxxx2")
}

type Intf interface {
	callf1(num int) int
	callf2(num int) int
}
type Ttf struct {
}

func (t Ttf) callf1(num int) int {
	fmt.Println("Ttf callf1", num)
	return num + 1
}
func (t Ttf) callf2(num int) int {
	fmt.Println("Ttf callf2", num)
	return num + 1
}

func test_interface() {
	ins := new(Ttf)
	ins.callf1(100)
	ins.callf2(200)

	var iii Intf
	iii = ins
	iii.callf1(300)
	iii.callf2(400)
}

func main() {
	//log_zap()
	test_interface()
	for {
		time.Sleep(2 * time.Second)
	}
}
