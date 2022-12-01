package main

import (
	"fmt"
	"os"
)

func isFileExsit(filename string) bool {
	_, err := os.Stat(filename)
	//fmt.Println("isFileExsit os.Stat=", err)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func main() {
	content := "11111122222\n"
	filename := "./out.txt"

	var err error

	//
	var file *os.File
	if isFileExsit(filename) {
		file, err = os.OpenFile(filename, os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("main open file err", err.Error())
			panic(err)
		}
	} else {
		file, err = os.Create(filename)
		if err != nil {
			fmt.Println("main create file err", err.Error())
			panic(err)
		}
	}
	defer func() {
		file.Close()
	}()
	var n int
	//n, err = io.WriteString(file, content)
	n, err = file.WriteString(content)
	fmt.Println("WriteString 写入n字节, n=", n, "err=", err)

	//
	//err = ioutil.WriteFile(filename, []byte(content), 0666)
	//fmt.Println("WriteFile err=", err)
}
