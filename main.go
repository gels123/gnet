package main

import (
	"fmt"
	"gnet/lib/logzap"
	"gnet/lib/utils"
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

func main() {

	fmt.Println(utils.GetExeDir())
	fmt.Println(utils.GetCurDir())

	//logOut := filepath.Join(utils.GetCurDir(), conf.LogsConf.FileDir, conf.LogsConf.FileName)
	//fmt.Println("===================df===", logOut)
	//fmt.Println("===================xxx===", conf.LogsConf.FileDir[0] == '.')

	////logzap.Infof("Failed to fetch URL: %s", "xxxx1")
	////logzap.Errorf("Failed to fetch URL: %s", "xxxx2")
	//
	//time.AfterFunc(time.Second*10, func() {
	//	fmt.Println("==============sdfadfadfa===============")
	//})

	//go aa()

	fmt.Println("=====================00=", time.Now().UTC().Unix())
	fmt.Println("=====================00=", time.Now().Unix())

	for {
		logzap.Debugw("========sdfadf===", "a=", 100)
		logzap.Debugf("========sdfadf===", "a=", 100)
		time.Sleep(time.Second)
	}
}
