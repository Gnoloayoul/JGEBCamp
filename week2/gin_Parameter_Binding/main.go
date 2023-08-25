package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	type UserInfo struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	r.LoadHTMLFiles("./index.html")
	r.GET("/user", func(c *gin.Context) {
		// 旧做法
		//username := c.Query("username")
		//password := c.Query("password")
		//u := UserInfo{
		//	Username: username,
		//	Password: password,
		//}

		// 新做法
		var u UserInfo
		err := c.ShouldBind(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			fmt.Printf("%#v\n", u)
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		}
	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/form", func(c *gin.Context) {
		var u UserInfo
		err := c.ShouldBind(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			fmt.Printf("%#v\n", u)
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		}
	})

	r.POST("/json", func(c *gin.Context) {
		var u UserInfo
		err := c.ShouldBind(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			fmt.Printf("%#v\n", u)
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		}
	})

	r.Run(":9090")
}
