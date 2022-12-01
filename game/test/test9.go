package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("./index.html")
	if err == nil {
		w.Write(content)
	} else {
		fmt.Fprintf(w, "welcome!!")
	}
}

func main() {
	//
	http.HandleFunc("/", index)
	err := http.ListenAndServe("0.0.0.0:8888", nil)
	if err != nil {
		fmt.Println("http.ListenAndServe err=", err)
		panic(err)
	}

}
