package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func sayhello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello golang!",
	})

}


func main() {
	// 准备路由引擎
	r := gin.Default()

	// 指定用户访问“/hello”时，执行sayhello函数
	r.GET("/hello", sayhello)

	//r.GET("/book",)
	//r.GET("/create_book",)
	//r.GET("/updata_book",)
	//r.GET("/delete_book",)

	// 推荐写法， Restful风格，需要用postman配合调试
	r.GET("/book", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"method": "Get",
		})
	})
	r.POST("/book", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"method": "POST",
		})
	})
	r.PUT("/book", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"method": "PUT",
		})
	})
	r.DELETE("/book", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"method": "DELETE",
		})
	})


	// 启动,默认是8080，要别的口，自己指定
	r.Run()
}