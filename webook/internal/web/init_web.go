package web

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	registerRoutes(server)
	return server
}

func registerRoutes(server *gin.Engine) {
	u := &UserHandler{}
	server.POST("/users/signup", u.SignUp)
	server.POST("/users/login", u.Login)
	server.POST("/users/edit", u.Edit)
	server.GET("/users/profile", u.Profile)
}

