package main

import (
	"errors"
	"fmt"
)

func fadd(num1 int, num2 int) int {
	return num1 + num2
}

func fcall(method string, nums ...int) {
	for k, v := range nums {
		fmt.Println(method, k, v)
	}
}

func walkSlice(slice []string, f func(k int, v string) bool) {
	for k, v := range slice {
		if !f(k, v) {
			break
		}
	}
	defer func() {
		fmt.Println("walkSlice end")
	}()
	//panic("crash")
}

func fdiv(a int, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("b == 0")
	}
	return a / b, nil
}

func main() {
	//函数指针
	var fptr func(int, int) int
	fptr = fadd
	fmt.Println("test3  sum=", fptr(100, 200))

	//
	fcall("get", 1, 2, 3)

	//
	//walkSlice([]string{"aa", "bb", "cc"}, func(k int, v string) bool {
	//	fmt.Println("walkSlice k=", k, "v=", v)
	//	if v == "bb" {
	//		return false
	//	} else {
	//		return true
	//	}
	//})

	//
	result, err := fdiv(100, 0)
	fmt.Println("div result=", result, "err=", err)
}
