package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func f1(w http.ResponseWriter, r *http.Request) {
	// 定义一个函数 kua
	k := func(name string) (string, error) {
		return name + "万人迷", nil
	}
	// 定义模板
	// 解析模板
	t, err := template.New("f.tmpl").Funcs(template.FuncMap{
		"kua": k,
	}).ParseFiles("./f.tmpl")
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	name := "首席"
	err = t.Execute(w, name)
	if err != nil {
		fmt.Printf("render template failed, err: %v\n", err)
		return
	}
}

func demo1(w http.ResponseWriter, r *http.Request) {
	// 定义模板
	// 解析模板
	t, err := template.ParseFiles("./t.tmpl", "./u1.tmpl")
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	name := "首席"
	err = t.Execute(w, name)
	if err != nil {
		fmt.Printf("render template failed, err: %v\n", err)
		return
	}
}

func main() {
	http.HandleFunc("/", f1)
	http.HandleFunc("/tmplDemo", demo1)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("http server failed, err: %v\n", err)
		return
	}
}