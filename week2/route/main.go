package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 路由和路由组

func main() {
	r := gin.Default()

	// 获取
	r.GET("/index", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"method": "GET",
		})
	})

	// 创建
	r.POST("/index", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"method": "POST",
		})
	})

	// 改 更新
	r.PUT("/index", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"method": "PUT",
		})
	})

	// 删
	r.DELETE("/index", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"method": "DELETE",
		})
	})

	// NoRoute
	r.NoRoute(func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"method": "404",
		})
	})

	// 视频的首页与详情页
	//r.GET("./video/index", func(c *gin.Context){
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "/video/index",
	//	})
	//})
	// 将公用的前缀提取出来，形成组
	videoGroup := r.Group("./video")
	{
		videoGroup.GET("./index", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "/video/index",
			})
		})
	}


	// 商场的首页与详情页
	r.GET("./shop/index", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "/shop/index",
		})
	})

	r.Run(":9090")
}