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

func TestFailoverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		wantErr error
	}{
		{
			name: "成功发送",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				s1 := smsmocks.NewMockService(ctrl)
				s2 := smsmocks.NewMockService(ctrl)
				s3 := smsmocks.NewMockService(ctrl)
				s := []sms.Service{s1, s2, s3}
				s1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("something wrong"))
				s2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return s
			},

			wantErr: nil,
		},
		{
			name: "全部失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				s1 := smsmocks.NewMockService(ctrl)
				s2 := smsmocks.NewMockService(ctrl)
				s3 := smsmocks.NewMockService(ctrl)
				s := []sms.Service{s1, s2, s3}
				s1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("something wrong"))
				s2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("something wrong"))
				s3.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("something wrong"))
				return s
			},

			wantErr: errors.New("全部服务商都失败"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ss := NewFailoverSMSService(tc.mock(ctrl))
			err := ss.Send(context.Background(), "mytpl", []string{"123"}, "135xxxxxxxx")

			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestFailoverSMSService_SendV1(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service
		idx  uint64

		wantErr error
		wantIdx uint64
	}{
		{
			name: "发生换源的正常发送",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				s1 := smsmocks.NewMockService(ctrl)

				s := []sms.Service{s1}
				s1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return s
			},
			idx: uint64(0),

			wantErr: nil,
			wantIdx: uint64(1),
		},
		{
			name: "全换过一遍，都失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				s0 := smsmocks.NewMockService(ctrl)
				s1 := smsmocks.NewMockService(ctrl)
				s2 := smsmocks.NewMockService(ctrl)
				s := []sms.Service{s0, s1, s2}
				s0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))
				s1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("换过源了，还是发送失败"))
				s2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("继续换，仍然失败"))
				return s
			},
			idx: uint64(0),

			wantErr: errors.New("全部服务商都失败"),
			wantIdx: uint64(1),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ss := NewFailoverSMSService(tc.mock(ctrl))
			ss.idx = tc.idx
			err := ss.SendV1(context.Background(), "mytpl", []string{"123"}, "135xxxxxxxx")

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantIdx, ss.idx)
		})
	}
}
