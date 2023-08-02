package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	regexp "github.com/dlclark/regexp2"
)

// UserHandler
// 与用户有关的路由
type UserHandler struct {
	emailExp *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		emailRegexPatten = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPatten = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPatten, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPatten, regexp.None)
	return &UserHandler{
		emailExp: emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	{
		ug.POST("/signup", u.SignUp)
		ug.POST("/edit", u.Edit)
		ug.GET("/profile", u.Profile)
		ug.POST("/delete", u.Delete)
		//// 临时signup.HTML用的
		//ug.GET("/index", u.Index)
	}

}

func (u *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		// 临时signup.HTML用的
		//Email string `form:"email"`
		//ConfirmPassword string `form:"confirmPassword"`
		//Password string `form:"password"`

		Email string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password string `json:"password"`
	}

	var req SignUpReq
	if err := c.ShouldBind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.String(http.StatusOK, "输入的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		c.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}

	c.String(http.StatusOK, "注册成功")
	fmt.Printf("%#v", req)
	// 接下来是数据库操作
}

//// 进入临时signup.HTML用的
//func (u *UserHandler) Index(c *gin.Context) {
//	 c.HTML(http.StatusOK, "signup.html", nil)
//}

func (u *UserHandler) Edit(c *gin.Context) {

}

func (u *UserHandler) Delete(c *gin.Context) {

}

func (u *UserHandler) Profile(c *gin.Context) {

}