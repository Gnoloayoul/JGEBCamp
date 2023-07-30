package main

import (
	"github.com/gin-go


	nic/gin"
	"html/template"
	"net/http"
)

func main() {
	r := gin.Default()
	// 加载静态文件
	r.Static("/xxx", "./statics")
	// gin框架中给模板追加自定义函数
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML{
			return template.HTML(str)
		},
	})
	// 解析模板
	//r.LoadHTMLFiles("templates/index.tmpl", "templates/user.tmpl")
	r.LoadHTMLGlob("templates/**/*")
	// 渲染模板
	r.GET("/post/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post/index.html", gin.H{
			"title": "github.com",
		})
	})

	r.GET("/user/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user/index.html", gin.H{
			//"title": "github.com",
			"title": "<a href='https://github.com'>首席的博客</a>",
		})
	})

	// 启动server引擎
	r.Run(":9000")
}









