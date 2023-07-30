package main

// from表单提交与获取参数
import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 一次请求，对应一个响应
// 一次请求，对应一个响应
// 一次请求，对应一个响应


func main() {
	r := gin.Default()
	r.LoadHTMLFiles("./login.html", "./index.html")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", func(c *gin.Context) {

		// 获取form表单提交的数据
		// 方法1：
		//username := c.PostForm("username")
		//password := c.PostForm("password")

		// 方法2：
		// 类似GET的DefaultQuery
		// 但是触发分支，不是不填（就返回空）
		// 而是就在DefaultPostForm里填入不存在的key
		//username := c.DefaultPostForm("username", "somebody")
		//password := c.DefaultPostForm("password", "????????")

		// 方法3：
		// 类似GET的GetQuery
		// 但是触发分支，不是不填（就返回空）
		// 而是就在DefaultPostForm里填入不存在的key
		// 这点和DefaultPostForm一样
		username, nok := c.GetPostForm("username")
		password, pok := c.GetPostForm("password")
		// 找不到名为username的key
		if !nok {
			username = "somebody"
		}
		// 找不到名为password的key
		if !pok {
			password = "*********"
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Name": username,
			"Password": password,
		})
	})


	r.Run(":9090")
}