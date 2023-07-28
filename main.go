package main

import (
	"fmt"
	"gnet/game/conf"
	"gnet/lib/log"
	"time"
)

func init() {
	fmt.Println("==== game main init begin ====")
	log.Init(conf.LogsConf.FileName, conf.LogsConf.FileName, conf.LogsConf.FileLevel, conf.LogsConf.ShellLevel, conf.LogsConf.MaxLine, conf.LogsConf.BufSize)
	fmt.Println("==== game main init end ====")
}

func main() {
	log.Info("==game main start==")
	log.Info("==game main start==")

	for {
		time.Sleep(10000000)
	}
}
