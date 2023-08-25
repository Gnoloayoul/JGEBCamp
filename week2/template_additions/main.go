package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	// 定义模板
	// 解析模板(在这里宣告要使用新的默认标识符)
	t, err := template.New("index.tmpl").Delims("{[", "]}").ParseFiles("./index.tmpl")
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	n := "灰鸦指挥官"
	err = t.Execute(w, n)
	if err != nil {
		fmt.Printf("render template failed, err: %v\n", err)
		return
	}
}

func xss(w http.ResponseWriter, r *http.Request) {
	// 定义模板
	// 解析模板
	t, err := template.ParseFiles("./xss.tmpl")
	if err != nil {
		fmt.Printf("Parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	// 演示html模板的安全能力，能把有恶意的语句转义掉
	src := "<script>alert(123);</script>"
	err = t.Execute(w, src)
	if err != nil {
		fmt.Printf("render template failed, err: %v\n", err)
		return
	}
}

func main() {
	http.HandleFunc("/index", index)
	http.HandleFunc("/xss", xss)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("http server failed, err: %v\n", err)
		return
	}
}
