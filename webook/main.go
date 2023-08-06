package main

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	u := web.NewUserHandler()
	u.RegisterRoutes(server)

	//// 临时用的signup页面
	//server.LoadHTMLFiles("../webook-fe/signup.html")

	server.Run(":8080")
}
