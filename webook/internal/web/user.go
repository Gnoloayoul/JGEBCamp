package web

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler
// 与用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPatten    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPatten = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPatten, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPatten, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	{
		ug.POST("/signup", u.SignUp)
		ug.POST("/edit", u.Edit)
		ug.GET("/profile", u.Profile)
		ug.POST("/login", u.Login)
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

		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	if err := c.Bind(&req); err != nil {
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

	err = u.svc.SignUp(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		c.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统异常")
		return
	}

	c.String(http.StatusOK, "注册成功")
}

//// 进入临时signup.HTML用的
//func (u *UserHandler) Index(c *gin.Context) {
//	 c.HTML(http.StatusOK, "signup.html", nil)
//}

func (u *UserHandler) Edit(c *gin.Context) {
	type InfoReq struct {
		Email    string `json:"email"`
		NickName           string `json:"nickname"`
		Birthday string `json:"birthday"`
		Info        string `json:"info"`
	}

	var req InfoReq
	if err := c.Bind(&req); err != nil {
		return
	}

	// 取得拿到userID
	sess := sessions.Default(c)
	id := sess.Get("userId").(int64)
	_, err := u.svc.Edit(c, domain.User{
		Id: id,
		NickName: req.NickName,
		Birthday: req.Birthday,
		Info: req.Info,
	})
	//if err == service.ErrInvalidUserOrPassword {
	//	c.String(http.StatusOK, "用户名或者密码不对, 不是当前的用户信息")
	//	return
	//}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	c.String(http.StatusOK, "补充信息修改成功")
}

func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.String(http.StatusOK, "用户名或者密码不对")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	// step 2
	// login success
	// set session
	sess := sessions.Default(c)
	// you can set the value what you want
	sess.Set("userId", user.Id)

	sess.Options(sessions.Options{
		// sess(在cookie里)保存多久？
		//
		MaxAge: 60,
	})

	sess.Save()
	c.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Profile(c *gin.Context) {
	// 取得拿到userID
	sess := sessions.Default(c)
	id := sess.Get("userId").(int64)

	userinfo, err := u.svc.Profile(c, domain.User{
		Id: id,
	})

	if err != nil {
		c.String(http.StatusOK, "系统错误, 用户信息404")
		return
	}

	c.JSON(http.StatusOK, userinfo)
}
