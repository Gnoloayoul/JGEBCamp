package web

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
	userSvc service.UserService
	jwtHandler
	stateKey []byte
	cfg Config
}

type Config struct {
	Secure bool
}

func NewWechatHandler(svc wechat.Service,
	userSvc service.UserService, cfg Config) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
		userSvc: userSvc,
		stateKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf1"),
		cfg: cfg,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat/")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.CallBack)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	// 要吧我的 state 存好
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg: "构造扫描登录 URL 失败",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, StateClaim{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			//  过期时间, 预期中一个用户完成登录的时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 3)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg: "系统错误",
		})
	}
	ctx.SetCookie("jwt-state", tokenStr,
		180, "/oauth2/wechat/callback",
		"", h.cfg.Secure, true)
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})


}

func (h *OAuth2WechatHandler) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	// 校验我之前的 state
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		// 确认有人攻击
		// 做好监控
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg: "系统错误",

		})
		return
	}
	info, err := h.svc.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg: "系统错误",

		})
		return
	}

	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg: "系统错误",

		})
		return
	}
	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg: "系统错误",

		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

type StateClaim struct {
	State string
	jwt.RegisteredClaims
}