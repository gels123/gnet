package main

import (
	"fmt"
	"gnet/lib/helper"
	"gnet/lib/loggerbak"
	"io/ioutil"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"gnet/lib/utils"
	"gnet/lib/uuid"
)

// 获取指定目录下的所有文件,包含子目录下的文件
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

type Student2 struct {
	name string
	age  int
}

func main() {
	//GetAllFiles("./test")
	fmt.Println("============", utils.GetTimeFormat(), utils.GetTimeW1H0())
	fmt.Println("=========sdfadf==", string(debug.Stack()))
	fmt.Println("=========sdfadf2==", helper.GetStack())
	fmt.Println("uuid1=", uuid.GenUUID(), "uuid2=", uuid.GenUUID(), uuid.GenFixedUUID([]byte("1")), uuid.GenFixedUUID([]byte("1")))

	logsimple.Init("test", "game", logsimple.DEBUG_LEVEL, logsimple.DEBUG_LEVEL, 10000, 1000)
	s := &Student2{"yyyyy", 100}
	logsimple.Error("hahaha %v, %v", 2, s)

	for {
		logsimple.Warn("loop log warn %v", 100)
		time.Sleep(2 * time.Second)
	}
}
