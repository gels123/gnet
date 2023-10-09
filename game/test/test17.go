package main

import (
	"fmt"
)

type itest interface {
	Say1()
	Say2()
}

type Person struct {
	son itest
}

func newPerson() *Person {
	return &Person{}
}

func (p *Person) setSuper(super itest) {
	p.son = super
}

func (p *Person) Say1() {
	fmt.Println("=Person Say1")
	p.Say2()
}

func (p *Person) Say2() {
	if p.son != nil {
		p.son.Say2()
		return
	}
	fmt.Println("=Person Say2")
}

type Doctor struct {
	base *Person
}

func newDoctor() *Doctor {
	d := &Doctor{base: newPerson()}
	//d.base.setSuper(d)
	return d
}

func (d *Doctor) Say1() {
	fmt.Println("=Doctor Say1")
	d.Say2()
	fmt.Println("==================xxxx")
	d.base.Say1()
}

func (d *Doctor) Say2() {
	fmt.Println("=Doctor Say2")
}

func main() {
	var a int = 100
	fmt.Println("==============", a)

	var t itest = newDoctor()
	t.Say1()

	//
	//ch := make(chan int)
	//go func() {
	//	time.Sleep(10 * time.Second)
	//	ch <- 1
	//}()
	//
	//<-ch
}
