package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type User struct {
	Name string
	Gender string
	Age int

}

func sayhello(w http.ResponseWriter, r *http.Request) {
	// 定义模板
	// 解析模板
	t, err := template.ParseFiles("./hello.tmpl")
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	// 想让模板文件能读到，字段首字母都得大写

	// 例子：结构体
	u1 := User{
		Name: "首席",
		Gender: "男",
		Age: 18,
	}

	// 例子：map
	m1 := map[string]interface{}{
		"Name": "首席",
		"Gender": "男",
		"Age": 18,
	}

	hobbylist := []string{
		"篮球",
		"足球",
		"双色球",
	}

	t.Execute(w, map[string]interface{}{
		"u1": u1,
		"m1": m1,
		"hobby": hobbylist,
	})
}

func main() {
	http.HandleFunc("/", sayhello)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("http server failed, err: %v\n", err)
		return
	}
}
