package main

import (
	"fmt"
	"gowork/lib/timer"
)

func main() {
	ts := timer.NewTimerSchedule()
	ts.Start()

	ts.Schedule(100, 10, func(dt int) {
		fmt.Println("time1", dt)
	})

	ch := make(chan int)
	<-ch
}
