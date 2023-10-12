package ratelimit

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ratelimit"
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
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rss := NewRatelimitSMSServiceV1(tc.mock(ctrl))
			err := rss.Send(context.Background(), "mytpl", []string{"123"}, "135xxxxxxxx")

			assert.Equal(t, tc.WantErr, err)
		})
	}
}
