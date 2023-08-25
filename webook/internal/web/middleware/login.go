package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

		// 怎么拿到sess的保质时间？
		// 不能使用time.time,这是单进程用的
		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now().UnixMilli()
		// 还没刷新（刚登录还么刷新）
		if updateTime == nil {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
			return
		}
		// 有
		updateTimeVal, _ := updateTime.(int64)

		if now-updateTimeVal > 60*1000 {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
			return
		}
	}
}

var IgnorePaths []string

func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不用登陆校验的
		for _, path := range IgnorePaths {
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

func CheckLoginV1(paths []string,
	abc int,
	bac int64,
	asdsd string) gin.HandlerFunc {
	if len(paths) == 0 {
		paths = []string{}
	}
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 不需要登录校验的
		//if ctx.Request.URL.Path == "/users/login" ||
		//	ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
