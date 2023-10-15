package failover

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	smsmocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

//func TestFailoverSMSService_Send(t *testing.T) {
//	testCases := []struct{
//		name string
//		mock func(ctrl *gomock.Controller) []sms.Service
//
//		wantErr error
//	}{
//		{
//			name: "成功发送",
//			mock: func(ctrl *gomock.Controller) []sms.Service {
//				 s1 := smsmocks.NewMockService(ctrl)
//				 s2 := smsmocks.NewMockService(ctrl)
//				 s3 := smsmocks.NewMockService(ctrl)
//				 s := []sms.Service{s1, s2, s3}
//				 s1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
//					 Return(errors.New("something wrong"))
//				 s2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
//					 Return(nil)
//				 return s
//			},
//
//			wantErr: nil,
//		},
//		{
//			name: "全部失败",
//			mock: func(ctrl *gomock.Controller) []sms.Service {
//				s1 := smsmocks.NewMockService(ctrl)
//				s2 := smsmocks.NewMockService(ctrl)
//				s3 := smsmocks.NewMockService(ctrl)
//				s := []sms.Service{s1, s2, s3}
//				s1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
//					Return(errors.New("something wrong"))
//				s2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
//					Return(errors.New("something wrong"))
//				s3.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
//					Return(errors.New("something wrong"))
//				return s
//			},
//
//			wantErr: errors.New("全部服务商都失败"),
//		},
//
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T){
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			ss := NewFailoverSMSService(tc.mock(ctrl))
//			err := ss.Send(context.Background(), "mytpl", []string{"123"}, "135xxxxxxxx")
//
//			assert.Equal(t, tc.wantErr, err)
//		})
//	}
//}

func TestFailoverSMSService_Send(t *testing.T) {
	testCases := []struct{
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		wantErr error
	}{
		{
			name: "全正常发送",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				s1 := smsmocks.NewMockService(ctrl)
				s2 := smsmocks.NewMockService(ctrl)
				s3 := smsmocks.NewMockService(ctrl)
				s1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				s2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				s3.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{s1, s2, s3}
			},
		},
		{
			name: "发生换源的正常发送",
			mock: func(ctrl *gomock.Controller) []sms.Service {

				s2 := smsmocks.NewMockService(ctrl)

				s := []sms.Service{s2}
				s2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return s
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ss := NewFailoverSMSService(tc.mock(ctrl))
			err := ss.Send(context.Background(), "mytpl", []string{"123"}, "135xxxxxxxx")

			assert.Equal(t, tc.wantErr, err)
		})
	}
}