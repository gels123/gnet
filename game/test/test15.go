package main

import (
	"fmt"
	"time"
)

type sbase struct {
	num int
}

func (base *sbase) print() {
	fmt.Println("sbase print", base.num)
}

type mycls struct {
	sbase
}

func (my *mycls) print() {
	fmt.Println("mycls print", my.num)
}

func main() {
	ch := make(chan int, 10)
	go f15_1(ch, 100)
	go f15_1(ch, 200)

	my := &mycls{
		sbase{num: 999},
	}
	my.print()
	my.sbase.print()

	for {
		time.Sleep(1 * time.Second)
	}
}

func f15_1(ch chan int, num int) {
	select {
	case ch <- num:
		fmt.Println("f15_1 send", num)
		break
	default:
		fmt.Println("f15_1 default", num)
	}
}
