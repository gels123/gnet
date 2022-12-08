package main

import (
	"fmt"
	"time"
)

var reqch chan int = make(chan int, 10)
var rspch chan string = make(chan string, 10)

func main() {
	//ts := timer.NewTimerSchedule()
	//ts.Start()
	//
	//ts.Schedule(100, 10, func(dt int) {
	//	fmt.Println("time1", dt)
	//})
	//

	go workerFunc()
	go factory()

	cch := make(chan int)
	<-cch
}

func workerFunc() {
	fmt.Println("==workerFunc==")
	for {
		data := <-reqch
		ret := fmt.Sprintf("{code=%v, data=%v}", 0, data)
		rspch <- ret
		fmt.Println("workerFunc process data=", data)
	}
}

func factory() {
	num := 100
	go func() {
		for {
			rsp := <-rspch
			fmt.Println("rec rsp=", rsp)
		}
	}()
	for {
		num += 1
		reqch <- num
		fmt.Println("send req=", num)
		time.Sleep(2 * time.Second)
	}
}

func call(addr chan interface{}, data interface{}) {
	addr <- data
}
