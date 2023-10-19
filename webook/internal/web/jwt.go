package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"time"
)

type jwtHandler struct {
	// access_token key
	atKey []byte
	// refresh_token key
	rtKey []byte
}

func newJwtHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		// 可以公用，也可以换一个
		rtKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	}
}

func (u jwtHandler) setRefashJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:       uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(u.rtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (u jwtHandler) setLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := u.setJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = u.setRefashJWTToken(ctx, uid, ssid)
	return err
}

func (u jwtHandler) setJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		Ssid: ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(u.atKey)
	if err != nil {
		return err
	}
	ctx.Header("x-JWT-token", tokenStr)
	return nil
}

func ExtraJWTToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	// 等效写法： sege := strings.SplitN(tokenHeader, " ", 2)
	// 因为这个 tokenHeader 也就两段
	sege := strings.Split(tokenHeader, " ")
	if len(sege) != 2 {
		return ""
	}
	return sege[1]
}

type RefreshClaims struct {
	Uid int64
	Ssid string
	jwt.RegisteredClaims
}

type UserClaims struct {
	jwt.RegisteredClaims
	Ssid string
	// 声明自己要放入的 token 里面的数据
	Uid       int64
	UserAgent string
}
