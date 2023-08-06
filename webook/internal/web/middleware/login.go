package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}


func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不用登陆校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

var IgnorePaths []string

func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不用登陆校验的
		for _, path := range IgnorePaths {

		}
	}
}