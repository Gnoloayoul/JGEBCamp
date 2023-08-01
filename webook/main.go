package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	"net/http"
)

func main() {
	server := gin.Default()
	// 临时用的signup页面
	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.LoadHTMLFiles("../webook-fe/signup.html")
	server.GET("/user/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})
	server.Run(":8080")
}