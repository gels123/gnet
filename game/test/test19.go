package main

import (
	"fmt"
	"reflect"
)

type struct1901 struct {
	num int
}

func (s *struct1901) GetNum() int {
	return s.num
}

func (s *struct1901) SetNum(num int) {
	s.num = num
}

func main() {
	fmt.Println("=============sdfasdf")

	s := &struct1901{
		num: 11,
	}
	fmt.Println("===========s.num=", s.num)

	rs := reflect.ValueOf(s)
	fmt.Println("===========rs=", rs)

	// 调用SetName方法
	ret := rs.MethodByName("SetNum").Call([]reflect.Value{reflect.ValueOf(22)})
	fmt.Println("===========s.SetNum=", ret)

	// 调用GetName方法
	ptr := rs.MethodByName("GetNum")
	fmt.Println("===========MethodByName=", ptr.IsValid())
	vals := rs.MethodByName("GetNum").Call([]reflect.Value{})
	for _, v := range vals {
		fmt.Println("===========s.GetNum=", v.Interface().(int))
	}
}

// 反射生成结构体
func GenStruct(req interface{}) interface{} {
	typ := reflect.TypeOf(req)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return reflect.New(typ).Interface()
}
