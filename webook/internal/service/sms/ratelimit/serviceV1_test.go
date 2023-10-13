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

func TestRatelimitSMSServiceV1_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter)

		wantErr error
	}{
		{
			name: "正常发送",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				sms := smsmocks.NewMockService(ctrl)
				limit := ratelimitmocks.NewMockLimiter(ctrl)

				limit.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				sms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return sms, limit
			},

			wantErr: nil,
		},
		{
			name: "限流功能出问题",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				sms := smsmocks.NewMockService(ctrl)
				limit := ratelimitmocks.NewMockLimiter(ctrl)

				limit.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("something wrong"))

				return sms, limit
			},

			wantErr: fmt.Errorf("短信服务判断是否限流出现问题， %w", errors.New("something wrong")),
		},
		{
			name: "触发限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				sms := smsmocks.NewMockService(ctrl)
				limit := ratelimitmocks.NewMockLimiter(ctrl)

				limit.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return sms, limit
			},

			wantErr: errLimited,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rss := NewRatelimitSMSServiceV1(tc.mock(ctrl))
			err := rss.Send(context.Background(), "mytpl", []string{"123"}, "135xxxxxxxx")

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
