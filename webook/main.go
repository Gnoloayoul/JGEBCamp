package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web/"
)

func main() {
	server := gin.Default()
	u := &web.UserHandler{}
	u.RegisterRouters(server)
	server.Run(":8080")
}