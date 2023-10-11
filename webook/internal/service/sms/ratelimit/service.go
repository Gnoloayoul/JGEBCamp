package ratelimit

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ratelimit"
)

// 装饰器模式

type Service struct {
	svc sms.Service
	limiter ratelimit.Limiter
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// 装饰器添加位置...
	err := s.svc.Send(ctx, tpl, args, numbers...)
	// 装饰器添加位置...
	return err
}




