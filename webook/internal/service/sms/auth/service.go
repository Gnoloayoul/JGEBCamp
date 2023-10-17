package auth

import (
	"context"
	"errors"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key string
}

// Send
// 其中 tpl （或者是 biz ）， 必须是线下申请的一个代表业务方的 Token
// 业务方拿着 tpl 过来，在这里解析
// 解析成功了，会换上对应的 Token（tpl），继续完成业务
func (s *SMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// 开始权限校验
	var tc Claims
	token, err := jwt.ParseWithClaims(tpl, &tc, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token 不合法")
	}
	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string

}