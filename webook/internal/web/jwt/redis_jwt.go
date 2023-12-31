package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// AccessTokenKey 因为 JWT Key 不太可能变，所以可以直接写成常量
// 也可以考虑做成依赖注入
var ATKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm")
var RTKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixA")

type RedisJwtHandler struct {
	cmd redis.Cmdable
}

func NewRedisJwtHandler(cmd redis.Cmdable) Handler {
	return &RedisJwtHandler{
		cmd: cmd,
	}
}

// SetLoginToken 设置登录后的 token
func (h *RedisJwtHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(ctx, uid, ssid)
	return err
}

func (h *RedisJwtHandler) SetRefreshToken(ctx *gin.Context,
	uid int64, ssid string,
) error {
	rc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置为七天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	refreshTokenStr, err := refreshToken.SignedString(RTKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil
}

// ClearToken 清除 token
func (h *RedisJwtHandler) ClearToken(ctx *gin.Context) error {
	// 正常用户的这两个 token 都会被前端更新
	// 也就是说在登录校验里面，走不到 redis 那一步就返回了
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	// 这里不可能拿不到
	uc := ctx.MustGet("user").(UserClaims)
	return h.cmd.Set(ctx, h.key(uc.Ssid),
		"", time.Hour*24*7).Err()
}

func (h *RedisJwtHandler) CheckSession(ctx *gin.Context, ssid string) error {
	logout, err := h.cmd.Exists(ctx,
		fmt.Sprintf("users:Ssid:%s", ssid)).Result()
	if err != nil {
		return err
	}
	if logout > 0 {
		return errors.New("用户已经退出登录")
	}
	return nil
}

func (h *RedisJwtHandler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return ""
	}
	// SplitN 的意思是切割字符串，但是最多 N 段
	// 如果要是 N 为 0 或者负数，则是另外的含义，可以看它的文档
	authSegments := strings.SplitN(authCode, " ", 2)
	if len(authSegments) != 2 {
		// 格式不对
		return ""
	}
	return authSegments[1]
}

func (h *RedisJwtHandler) SetJWTToken(ctx *gin.Context, uid int64,
	ssid string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Id:        uid,
		Ssid:      ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 演示目的设置为一分钟过期
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			// 在压测的时候，要将过期时间设置更长一些
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	})
	tokenStr, err := token.SignedString(ATKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

//type RedisHandler struct {
//	cmd redis.Cmdable
//	// 长 token 的过期时间
//	rtExpiration time.Duration
//}
//
//func NewRedisHandler(cmd redis.Cmdable) Handler {
//	return &RedisHandler{
//		cmd:          cmd,
//		rtExpiration: time.Hour * 24 * 7,
//	}
//}

func (h *RedisJwtHandler) key(ssid string) string {
	return fmt.Sprintf("users:Ssid:%s", ssid)
}
