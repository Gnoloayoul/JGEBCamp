package failover

import (
	"context"
	"errors"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	smsmocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestTimeoutFailoverSMSService_Send(t *testing.T) {
	testCases := []struct{
		name string
		mock func(ctrl *gomock.Controller) []sms.Service
		threshold int32
		idx int32
		cnt int32

		wantErr error
		wantIdx int32
		wantCnt int32
	}{
		{
			name: "超时，但是没连续超时",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)
				return []sms.Service{svc0}
			},
			threshold: 3,
			wantErr: context.DeadlineExceeded,
			wantCnt: 1,
			wantIdx: 0,
		},
		{
			name: "触发了切换，切换之后成功了",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0, svc1}
			},
			threshold: 3,
			cnt: 3,

			wantCnt: 0,
			wantIdx: 1,
		},
		{
			name: "触发了切换，但是还是失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))
				return []sms.Service{svc0, svc1}
			},
			threshold: 3,
			cnt: 3,

			wantErr: errors.New("发送失败"),
			wantCnt: 0,
			wantIdx: 1,
		},
		{
			name: "触发了切换，但是还是超时",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)
				return []sms.Service{svc0, svc1}
			},
			threshold: 3,
			cnt: 3,

			wantErr: context.DeadlineExceeded,
			wantCnt: 1,
			wantIdx: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tss := NewTimeoutFailoverSMSService(tc.mock(ctrl), tc.threshold)
			tss.idx = tc.idx
			tss.cnt = tc.cnt

			err := tss.Send(context.Background(), "mytpl", []string{"123"}, "135xxxxxxxx")

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantIdx, tss.idx)
			assert.Equal(t, tc.wantCnt, tss.cnt)
		})
	}
}
