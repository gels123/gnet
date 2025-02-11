package main

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"time"
)

func init() {
}

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

func aa() {
	bb()
}

func bb() {
	cc()
}

func cc() {
	gg()
}

func gg() {
	fmt.Println("==========sdfasdfadsf===gg=", string(debug.Stack()))

	t1 := time.Now()
	t2 := t1.Add(time.Second * 5)
	time.AfterFunc(t2.Sub(t1), gg)
}

// func main() {
// 	fmt.Println("===main begin===" + utils.GetExeDir())

// 	// logOut := filepath.Join(utils.GetCurDir(), conf.LogsConf.FileDir, conf.LogsConf.FileName)
// 	// fmt.Println("===================df===", logOut)
// 	// fmt.Println("===================xxx===", conf.LogsConf.FileDir[0] == '.')

// 	// time.AfterFunc(time.Second*10, func() {
// 	// 	fmt.Println("==============sdfadfadfa===============")
// 	// })

// 	// logzap.SetSource("gelsxxx")
// 	// logzap.Panic("dffffffffffffffffff")
// 	// logzap.Infow("========sdfadf===", "num=", 100)
// 	// logzap.Error("========sdfadf===", zap.String("nnnn", "nihao"))
// 	// for {
// 	// 	//logzap.Debugw("=sdfadf=", "num=", 100)
// 	// 	//logzap.Debugf("========sdfadf===", "a=", 100)
// 	// 	time.Sleep(time.Second)
// 	// }

// 	// queue := lockfreequeue.NewQueue(1024 * 1024)
// 	// queue.Put(100)
// 	// queue.Put(200)
// 	// queue.Put(300)
// 	// queue.Put(400)
// 	// queue.Put(500)
// 	// queue.Put(600)
// 	// cell, _, _ := queue.Get()
// 	// var isok bool = true
// 	// fmt.Println("-------", cell, !isok)


// 	fmt.Println("====main end===")

// }
