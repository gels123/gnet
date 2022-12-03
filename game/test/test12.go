package main

import (
	"fmt"
	"gowork/engine/log"
	"io/ioutil"
	"strconv"
	"strings"

	//"gowork/engine/gwlog"
	"gowork/engine/utils"
	"gowork/engine/uuid"
	"time"
)

//获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return
	}
	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
		} else {
			// 过滤指定格式
			fname := fi.Name()
			ok := strings.HasPrefix(fname, "game-2022-12-03.log")
			if ok {
				n := strings.LastIndex(fname, ".")
				fname = fname[n+1:]
				ii, err := strconv.Atoi(fname)
				fmt.Println("=========sdfdsf==", fi.Name(), ii, err)
			}
		}
	}
}

func main() {
	//GetAllFiles("./test")
	fmt.Println("============", utils.GetTimeFormat(), utils.GetTimeW1H0())
	fmt.Println("uuid1=", uuid.GenUUID(), "uuid2=", uuid.GenUUID(), uuid.GenFixedUUID([]byte("1")), uuid.GenFixedUUID([]byte("1")))

	//gwlog.Debugf("this is a debug %d", 1)

	log.Init("test", "game", log.DEBUG_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	for {
		time.Sleep(time.Second)
		log.Error("hahaha %v, %v", 2, 3)
	}
}
