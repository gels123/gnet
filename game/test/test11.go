package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var ch chan int32
var rw sync.RWMutex
var val int32

func main() {
	ch = make(chan int32, 10)
	go write(ch)
	go read(ch)
	for {
		time.Sleep(10000000)
	}
}

func write(ch chan int32) {
	for {
		rw.Lock()
		num := rand.Int31n(1000)
		//ch <- num
		val = num
		fmt.Println("write num=", num)
		rw.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func read(ch chan int32) {
	for {
		rw.RLock()
		//num := <-ch
		num := val
		fmt.Println("read num=", num)
		rw.RUnlock()
		time.Sleep(1 * time.Second)
	}
}
