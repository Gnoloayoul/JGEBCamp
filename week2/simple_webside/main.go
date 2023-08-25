package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func sayhello(w http.ResponseWriter, r *http.Request) {
	outPut, _ := ioutil.ReadFile("./hello.txt")
	_, _ = fmt.Fprintln(w, string(outPut))
}

func main() {
	http.HandleFunc("/hello", sayhello)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Printf("http server failed, err:%v\n", err)
		return
	}
}
