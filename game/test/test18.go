package main

import (
	"fmt"
	"reflect"
)

func test18_f1() (i int) { // 返回200
	defer func() {
		i = 200
	}()
	defer func() {
		i = 300
	}()
	i = 100
	return 1
}

func test18_f2() int {
	var str string = "abc中文"
	fmt.Println("len(str) =", len(str))                 // 9=3+3*2,go默认使用utf-8编码
	fmt.Println("len([]rune(str)) =", len([]rune(str))) // 5=3+1+1
	return 1
}

func test18_f3() {
	type stest struct {
		name string `json:name`
		age  int    `json:age`
	}
	t := &stest{"dsfadf", 100}
	filed, ok := reflect.TypeOf(t).Elem().FieldByName("name")
	if !ok {
		panic("FieldByName find filed failed")
	}
	fmt.Println("test filed=", filed.Type, filed.Name, string(filed.Tag))
}

func test18_f4(i *int) {
	*i = 200 // 会改变外部i的值
}

func test18_f5(slice []int) {
	slice = append(slice, 100) // append可能会扩容, 导致地址发生变化, 外部slice不会被修改
}

func test18_f6(slice []int) {
	slice[0] = 123 //地址一致，外部slice会被修改
}

func main() {
	fmt.Println("===============sdfadf1====", test18_f1())
	test18_f2()
	test18_f3()

	i := 100
	test18_f4(&i)
	fmt.Println("--------i=", i)

	slice := make([]int, 1)
	slice[0] = 0
	test18_f5(slice)
	test18_f6(slice)
	fmt.Println("------------slice=", slice)
}
