package gRPC

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	testCases := []struct {
		name  string
		reqId int64
	}{
		{
			name:  "user:123",
			reqId: 123,
		},
		{
			name:  "user:456",
			reqId: 456,
		},
		{
			name:  "user:444",
			reqId: 444,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建gRPC连接
			cc, err := grpc.Dial(":8090",
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			client := NewUserSvcClient(cc)
			_, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			resp, err := client.GetById(context.Background(), &GetByIdRep{Id: tc.reqId})
			if err != nil {
				// 检查错误类型
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("Failed to get gRPC status: %v", err)
				}

				// 检查错误代码
				if st.Code() == codes.NotFound {
					t.Log(tc.name, "Error:", st.Message())
				} else {
					t.Fatalf("Unexpected error: %v", err)
				}
			} else {
				t.Log(tc.name, resp.User)
			}
		})
	}
}
