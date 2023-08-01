package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler
// 与用户有关的路由
type UserHandler struct {

}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	{
		ug.POST("/signup", u.SignUp)
		ug.POST("/edit", u.Edit)
		ug.GET("/profile", u.Profile)
		ug.POST("/delete", u.Delete)
	}

}

func (u *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Email string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password string `json:"password"`
	}

	var req SignUpReq
	if err := c.Bind(&req); err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"mes": "注册成功",
	})
	fmt.Printf("%v", req)
}

func (u *UserHandler) Edit(c *gin.Context) {

}

func (u *UserHandler) Delete(c *gin.Context) {

}

func (u *UserHandler) Profile(c *gin.Context) {

}