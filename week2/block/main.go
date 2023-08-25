package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	// 定义模板
	// 解析模板
	t, err := template.ParseFiles("./t.tmpl")
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	msg := "空花魅魔"
	err = t.Execute(w, msg)
	if err != nil {
		fmt.Printf("render template failed, err: %v\n", err)
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	// 定义模板
	// 解析模板
	t, err := template.ParseFiles("./home.tmpl")
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	msg := "空花魅魔"
	err = t.Execute(w, msg)
	if err != nil {
		fmt.Printf("render template failed, err: %v\n", err)
		return
	}
}

func index2(w http.ResponseWriter, r *http.Request) {
	// 定义模板
	// 解析模板(模板继承)
	t, err := template.ParseFiles("./templates/base.tmpl", "./templates/index.tmpl")
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	msg := "空花魅魔"
	err = t.ExecuteTemplate(w, "index.tmpl", msg)
	if err != nil {
		fmt.Printf("render template failed, err: %v\n", err)
		return
	}
}

func home2(w http.ResponseWriter, r *http.Request) {
	// 定义模板
	// 解析模板(模板继承)
	t, err := template.ParseFiles("./templates/base.tmpl", "./templates/home.tmpl")
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	msg := "空花魅魔"
	err = t.ExecuteTemplate(w, "home.tmpl", msg)
	if err != nil {
		fmt.Printf("render template failed, err: %v\n", err)
		return
	}
}

func main() {
	http.HandleFunc("/index", index)
	http.HandleFunc("/home", home)
	http.HandleFunc("/index2", index2)
	http.HandleFunc("/home2", home2)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("http server failed, err: %v\n", err)
		return
	}
}
