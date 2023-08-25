package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 重定向

func main() {
	r := gin.Default()

	r.GET("/index", func(c *gin.Context) {
		//c.JSON(http.StatusOK, gin.H{
		//	"status": "ok",
		//})

		// 把请求转给别人
		c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})

	r.GET("/a", func(c *gin.Context) {
		// 跳转到/b的路由函数

		// 直接修改路由地址
		c.Request.URL.Path = "/b"
		// 按照新路由开始处理
		r.HandleContext(c)
	})

	r.GET("/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "b",
		})
	})

	r.Run(":9090")
}
