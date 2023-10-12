package ratelimit

import (
	"context"
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ratelimit"
)

type RatelimitSMSServiceV1 struct {
	sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSServiceV1(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService {
		svc: svc,
		limiter: limiter,
	}
}

func (s *RatelimitSMSServiceV1) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流出现问题， %w", err)
	}
	if limited {
		return errLimited
	}
	err = s.Service.Send(ctx, tpl, args, numbers...)
	return err
}
