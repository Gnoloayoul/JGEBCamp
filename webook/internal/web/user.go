package web

import (
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	userIdKey           = "userId"
	bizLogin            = "login"
	emailRegexPatten    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPatten = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

// UserHandler
// 与用户有关的路由
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	jwtHandler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	emailExp := regexp.MustCompile(emailRegexPatten, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPatten, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		codeSvc:     codeSvc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	{
		ug.POST("/signup", u.SignUp)
		ug.POST("/edit", u.Edit)
		//ug.GET("/profile", u.Profile)
		ug.GET("/profile", u.ProfileJWT)
		//ug.POST("/login", u.Login)
		ug.POST("/login", u.LoginJWT)
		//// 临时signup.HTML用的
		//ug.GET("/index", u.Index)
		ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
		ug.POST("/login_sms/code/loginsms", u.LoginSMS)
	}
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 这边，可以加上各种校验
	ok, err := u.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}

	// 我这个手机号，会不会是一个新用户呢？
	// 这样子
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 这边要怎么办呢？
	// 从哪来？
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 是不是一个合法的手机号码
	// 考虑正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}

	err := u.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
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
		c.String(http.StatusOK, "系统异常")
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
		c.String(http.StatusOK, "系统异常")
		return
	}
	if !ok {
		c.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符、字母")
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
		NickName string `json:"nickname"`
		Birthday string `json:"birthday"`
		Info     string `json:"info"`
	}

	var req InfoReq
	if err := c.Bind(&req); err != nil {
		return
	}

	// 取得拿到userID
	sess := sessions.Default(c)
	id := sess.Get("userId").(int64)
	_, err := u.svc.Edit(c, domain.User{
		Id:       id,
		NickName: req.NickName,
		Birthday: req.Birthday,
		Info:     req.Info,
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
		Secure:   true,
		HttpOnly: true,
		MaxAge:   60,
	})

	sess.Save()
	c.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 步骤2
	// 在这里用 JWT 设置登录态
	// 生成一个 JWT token

	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Println(user)
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Profile(c *gin.Context) {
	// 取得拿到userID
	sess := sessions.Default(c)
	id := sess.Get("userId").(int64)

	userinfo, err := u.svc.Profile(c, id)

	if err != nil {
		c.String(http.StatusOK, "系统错误, 用户信息404")
		return
	}

	c.JSON(http.StatusOK, userinfo)
}

func (u *UserHandler) ProfileJWT(c *gin.Context) {
	//// 取得拿到userID
	//sess := sessions.Default(c)
	//id := sess.Get("userId").(int64)

	cla, ok := c.Get("claims")
	if !ok {
		// 假设 这里没有拿到 claims， 怎么办？
		c.String(http.StatusOK, "系统错误")
		return
	}
	// ok 代表是不是 *UserClaims
	claims, ok := cla.(*UserClaims)
	if !ok {
		// 监控这里
		c.String(http.StatusOK, "系统错误")
		return
	}

	fmt.Println(claims.Uid)

	c.String(http.StatusOK, "你的 profile")
}

