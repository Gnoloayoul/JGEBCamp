package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"

)

func main() {
	server := gin.Default()

	u := web.NewUserHandler()
	u.RegisterRoutes(server)

	//// 临时用的signup页面
	//server.LoadHTMLFiles("../webook-fe/signup.html")

	server.Run(":8080")
}