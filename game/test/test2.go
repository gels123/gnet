package main

import (
	"container/list"
	"fmt"
	"sort"
	"sync"
	"time"
)

func f1() {
	//
	var arr1 [5]int
	for i := 1; i < 5; i++ {
		arr1[i] = 100 + i
	}
	fmt.Println("====f1=== arr1=", arr1)

	//
	//var arr2 = []int{1, 2, 3, 4, 5, 6}
	arr2 := []int{1, 2, 3, 4, 5, 6}
	fmt.Println("====f1=== arr2=", len(arr2), cap(arr2), arr2[1:3])

	//10,10=len,cap
	arr3 := make([]int, 10, 10)
	for i := 0; i < cap(arr3); i++ {
		arr3[i] = 1000 + i
	}
	arr3 = append(arr3, 100, 200)
	sort.Ints(arr3)
	fmt.Println("=====f1==== arr3=", arr3, cap(arr3), len(arr3))

	//
	slice1 := make([]int, 3, 5)
	for i, _ := range slice1 {
		slice1[i] = 500 + i
	}
	fmt.Println("=====f1==== slice1=", slice1, cap(slice1), len(slice1))

	//map非线程安全
	//var map1 map[string]int
	//var map1 = make(map[string]int)
	map1 := make(map[string]int, 100) //size
	map1["lili"] = 100
	map1["lilei"] = 99
	map2 := make(map[int][]string, 5) //map key=int, map val=[]string
	map2[1] = make([]string, 1, 10)
	map2[1][0] = "ni"
	map2[1] = append(map2[1], "haohao")
	//delete(map2, 1)
	fmt.Println("=====f1==== map1=", map1, len(map1), "map2=", map2, len(map2))
	for k, v := range map2 {
		fmt.Println("=====f1==== map2 k=", k, "v=", v)
	}

	//
	var synmap sync.Map
	synmap.Store("aa", 101)
	synmap.Store("bb", 102)
	synmap.Store("cc", 103)
	synmap.Delete("cc")
	synmap.Range(func(k, v interface{}) bool {
		fmt.Println("=====f1 synmap k=", k, "v=", v)
		return true
	})

	//
	//var list1 list.List
	list1 := list.New()
	list1.PushBack("hello")
	list1.PushBack("mouto")
	fmt.Println("====f1 list1=", list1, list1.Len())
	for it := list1.Front(); it != nil; it = it.Next() {
		fmt.Println("====f1 list1 it=", it.Value)
	}

	//
	for {
		for i := 1; i < 10; i++ {
			if i == 9 {
				goto mgoto
			}
		}
	}
mgoto:
	fmt.Println("====f1 markgoto")
}

func f2() {
	stime := time.Now()
	fmt.Println("=========f2=======start, stime=", stime)
	slice1 := make([]int, 10, 10)
	for i := 0; i < len(slice1); i++ {
		slice1[i] = 100 + i
	}
	fmt.Println("f2 slice1=", slice1)
	for k, v := range slice1 {
		go f3(&k, &v) //错误, 应该传递值, 不能传递地址
	}
	costtime := time.Since(stime)
	fmt.Println("=========f2=======end, costtime=", costtime)
}

func f3(k *int, v *int) {
	time.Sleep(1)
	fmt.Println("f3 k=", *k, "v=", *v)
}

func main() {
	f1()
	//f2()
	//for {
	//	time.Sleep(1)
	//}
}
