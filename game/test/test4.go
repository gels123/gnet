package main

import (
	"fmt"
	"runtime"
)

func ProtectRun(f func(val interface{}), val interface{}) {
	defer func() {
		err := recover()
		switch err.(type) {
		case runtime.Error:
			fmt.Println("is runtime.Error", err)
			break
		default:
			fmt.Println("is not runtime.Error", err)
		}
	}()
	f(val)
}
func main() {
	fmt.Println("============111====")
	ProtectRun(func(val interface{}) {
		fmt.Println("ProtectRun enter=", val)
		total := 100
		if val != nil {
			v, ok := val.(int)
			if ok {
				total = total / v
			}
		}
		fmt.Println("ProtectRun end total=", total)
	}, 0)
	fmt.Println("============222====")
}
