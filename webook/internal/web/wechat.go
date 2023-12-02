package web

import (
	"errors"
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/oauth2/wechat"
	ijwt "github.com/Gnoloayoul/JGEBCamp/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
	stateKey []byte
	//cfg      WechatHandlerConfig
}

//type WechatHandlerConfig struct {
//	Secure bool
//}

func NewOAuth2WechatHandler(svc wechat.Service,
	userSvc service.UserService,
	jwtHdl ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf1"),
		//cfg:      cfg,
		Handler:  jwtHdl,
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
			Msg:  "构造扫描登录 URL 失败",
		})
		return
	}
	if err = h.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫描登录 URL 失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})

}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, StateClaim{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			//  过期时间, 预期中一个用户完成登录的时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 3)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}
	ctx.SetCookie("jwt-state", tokenStr,
		180, "/oauth2/wechat/callback",
		"", false, true)
	return nil
}

func (h *OAuth2WechatHandler) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}

	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	// 校验我之前的 state
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到 State 的 cookie, %w", err)
	}

	var sc StateClaim
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token 已经过期, %w", err)
	}

	if sc.State != state {
		return errors.New("state 不相等")
	}
	return nil
}

type StateClaim struct {
	State string
	jwt.RegisteredClaims
}
