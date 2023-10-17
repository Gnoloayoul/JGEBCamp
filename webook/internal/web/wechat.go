package web

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
	userSvc service.UserService
	jwtHandler
}

func NewWechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat/")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.CallBack)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg: "构造扫描登录 URL 失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})


}

func (h *OAuth2WechatHandler) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	info, err := h.svc.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg: "系统错误",

		})
		return
	}

	u, err := h.userSvc.FindOrCreateByWechat(ctx, info.OpenID)
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