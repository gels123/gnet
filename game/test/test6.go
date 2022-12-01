package main

import (
	"fmt"
	"log"
	"runtime"
)

type Student struct {
	name  string
	age   uint
	score uint
}

func NewStudent(name string, age uint, score uint) *Student {
	return &Student{
		name:  name,
		age:   age,
		score: score,
	}
}

func (s *Student) Print() {
	fmt.Println("Student Print=", s.name, s.age, s.score)
}

type GodStudent struct {
	Student
	reward uint
}

func NewGodStudent(name string, age uint, score uint, reward uint) *GodStudent {
	//ins := &GodStudent{}
	//ins.name, ins.age, ins.score, ins.reward = name, age, score, reward
	//return ins
	return &GodStudent{
		Student: Student{name: name, age: age, score: score},
		reward:  reward,
	}
}

func (s *GodStudent) Print() {
	fmt.Println("GodStudent Print=", s.name, s.age, s.score, s.reward)
}

func main() {
	//
	lili := Student{"lili", 20, 99}
	lilei := Student{"lilei", 20, 90}
	var xiaoming Student
	xiaoming.name, xiaoming.age, xiaoming.score = "xiaoming", 21, 80
	fmt.Println("==lili==", lili, "lilei==", lilei, "xiaoming=", xiaoming, "lihong==", Student{"lihong", 19, 60})

	//
	//lisa := new(Student)
	lisa := &Student{"lisa", 20, 90}
	//var lisa *Student = new(Student)
	lisa.name, lisa.age, lisa.score = "lisa", 20, 90
	fmt.Println("==lisa==", *lisa)

	//
	{
		var s *GodStudent = NewGodStudent("xiaohua", 20, 100, 10)
		runtime.SetFinalizer(s, func(s *GodStudent) {
			s.Print()
		})
		fmt.Println("====godStu====", *s)
		s.Print()
		s.Student.Print()
	}

	//
	fmt.Println("===GC==")
	runtime.GC()
	//runtime.Gosched()
	log.Println("==end==")
}
