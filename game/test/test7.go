package main

import (
	"fmt"
	"sort"
)

type MyWriter interface {
	Write(data string) error
}

type MySock struct {
}

func (s *MySock) Write(data string) error {
	fmt.Println("MySock Write=", data)
	return nil
}

type Player struct {
	pid  int
	name string
}

func test7_f1(v interface{}) {
	switch v.(type) {
	default:
		fmt.Println("testf1===", v)
	}
}

func main() {
	var mysock *MySock = &MySock{}
	var writer MyWriter = mysock
	writer.Write("hello world")

	playerList := []Player{{101, "lord101"}, {111, "lord111"}, {100, "lord100"}}
	sort.Slice(playerList, func(i, j int) bool {
		return playerList[i].pid < playerList[j].pid
	})
	fmt.Println("====playerList====", playerList)

	test7_f1(100)
}
