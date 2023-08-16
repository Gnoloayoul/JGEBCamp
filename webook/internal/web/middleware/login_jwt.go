package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

// LoginJWTMiddlewareBuilder
// 基于 JWT 的登陆校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不用登陆校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		// 用 JWT 做登陆校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 等效写法： sege := strings.SplitN(tokenHeader, " ", 2)
		// 因为这个 tokenHeader 也就两段
		sege := strings.Split(tokenHeader, " ")
		if len(sege) != 2 {
			// 有人捣乱
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 正式取得token
		tokenStr := sege[1]
		// 还原 jwt
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("h7oUXRzcGPyJbZJfq68iGChnzA0iJBfJ"), nil
		})
		if err != nil {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// err 为 nil， token 不为 nil
		if token == nil || !token.Valid {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}




