package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	smsmocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms/mocks"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ratelimit"
	ratelimitmocks "github.com/Gnoloayoul/JGEBCamp/webook/pkg/ratelimit/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRatelimitSMSService_Send(t *testing.T) {
	testCases := []struct{
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter)

		wantErr error
	}{
		{
			name: "没触发限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limiter := ratelimitmocks.NewMockLimiter(ctrl)

				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, limiter
			},

			wantErr: nil,
		},
		{
			name: "限流出现问题",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limiter := ratelimitmocks.NewMockLimiter(ctrl)

				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("something worng"))
				return svc, limiter
			},

			wantErr: fmt.Errorf("短信服务判断是否限流出现问题， %w", errors.New("something worng")),
		},
		{
			name: "触发限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limiter := ratelimitmocks.NewMockLimiter(ctrl)

				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return svc, limiter
			},

			wantErr: errors.New("触发限流"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rss := NewRatelimitSMSService(tc.mock(ctrl))
			err := rss.Send(context.Background(), "mytpl", []string{"123"}, "135xxxxxxxx")

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
