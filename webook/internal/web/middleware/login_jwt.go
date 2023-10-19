package middleware

import (
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
)

// LoginJWTMiddlewareBuilder
// 基于 JWT 的登陆校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
	cmd redis.Cmdable
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
		// 正式取得token(提取)
		tokenStr := web.ExtraJWTToken(ctx)
		// 这里要用指针，因为下面的 ParseWithClaims 就是会修改里面的数值
		claims := &web.UserClaims{}
		// 还原 jwt
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("h7oUXRzcGPyJbZJfq68iGChnzA0iJBfJ"), nil
		})
		if err != nil {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// err 为 nil， token 不为 nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 需要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		cnt, err := l.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid)).Result()
		if err != nil || cnt > 0  {
			// redis 出问题， 或者你已经退出登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("claims", claims)
		//ctx.Set("userId", claims.Uid)
	}
}
