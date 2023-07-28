package main

import (
	"fmt"
	"math"
)

func main() {
	var a, b int
	a = 100
	b = 200
	var num1, num2 *int
	num1 = &a
	num2 = &b

	var (
		c int
		d int
	)
	c = 300
	d = 300

	e, f := 500, 600
	e, f = f, e
	fmt.Println("==11111==", *num1, *num2, c, d, e, f)

	sum, _, _ := add(e, f)
	fmt.Println("sum =", sum)

	var fnum1 float64 = 100 / 200.0
	const constnum int32 = 999
	fmt.Println("==222222==fnum1=", fnum1, "constnum=", constnum)

	//复数
	com1, com2 := complex(1, 2), complex(2, 4)
	fmt.Println("complex=", com1+com2, com1*com2, math.Sin(60))

	//bool
	var isOk bool = false
	fmt.Println("bool=", isOk, !isOk, isOk == true, isOk == false)

	//字符串
	var str1 string = "hello world"
	str2 := `
		asdfad
		xxxxxx
	`
	fmt.Println("str1=", str1, str2, addstr(str1, str2), len(str2))

	//类型定义
	type myint = int32
	var my1 myint = 100
	fmt.Println("my1=", my1)

	//不同类型不同处理
	var ii int = 120
	var vi interface{} = ii
	switch vi.(type) {
	case string:
		n, err := vi.(string)
		fmt.Println("====1 vi to string", n, err)
		break
	case int:
		n, err := vi.(int)
		fmt.Println("====2 vi to int", n, err)
		break
	default:
		fmt.Println("====3 vi to default")
	}
}

func add(a int, b int) (int, int, int) {
	return a + b, a, b
}

func addstr(str1 string, str2 string) string {
	return str1 + str2
}
