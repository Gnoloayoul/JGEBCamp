package ratelimit

import (
	"context"
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ratelimit"
)

var errLimited = fmt.Errorf("触发限流")

// 装饰器模式

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc: svc,
		limiter: limiter,
	}
}

func (s *RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流出现问题， %w", err)
	}
	if limited {
		return errLimited
	}
	// 装饰器添加位置...
	err = s.svc.Send(ctx, tpl, args, numbers...)
	// 装饰器添加位置...
	return err
}
