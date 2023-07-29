package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func sayhello(w http.ResponseWriter, r *http.Request) {
	// 2.解析模板
	t, err := template.ParseFiles("./hello.tmpl") // 小心刻舟求剑
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 3.定义模板
	name := "空花魅魔 首席"
	err = t.Execute(w, name)
	if err != nil {
		fmt.Printf("render template failed, err: %v\n", err)
		return
	}
}

func main() {
	http.HandleFunc("/", sayhello)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("http server failed, err: %v\n", err)
		return
	}
}

// 这里最好go build
// 用go run的话，二进制文件都不知道跑到哪里，那么模板文件的路径就是错的了