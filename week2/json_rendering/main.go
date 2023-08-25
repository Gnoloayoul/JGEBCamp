package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/json", func(c *gin.Context) {
		// 方法1: 使用map
		//data := map[string]interface{} {
		//	"name": "首席",
		//	"message": "helloworld",
		//	"age": 10,
		//}

		// 使用Gin提供的快捷方式：用gin.H对map[string]interface{}封装
		//data := gin.H{"name":"首席", "message":"Hello", "age":"18"}

		// 方法2： 使用结构体
		type msg struct {
			Name    string `json:"name"`
			Message string
			Age     int
		}
		data := msg{
			"首席",
			"hello sec",
			44,
		}
		c.JSON(http.StatusOK, data)
	})

	r.Run(":9090")
}
