package main

// querystring

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/web", func(c *gin.Context) {
		// 获取浏览器那边发送请求携带的querystring参数
		//name := c.Query("query") // 方法1： 通过Query获取请求中的querystring参数
		//name := c.DefaultQuery("query", "somebody") // 方法2： 使用defaultQuery，取不到的就返回默认值somebody
		name, nok := c.GetQuery("query") // 方法3： 使用getquery，取不到返回false
		age, aok := c.GetQuery("age")
		if !nok {
			// 取不到
			name = "somebody"
		}
		if !aok {
			// 取不到
			age = "???"
		}
		c.JSON(http.StatusOK, gin.H{
			"name": name,
			"age": age,
		})

	})

	r.Run(":9090")
}