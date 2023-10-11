package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	"github.com/Gnoloayoul/JGEBCamp/webook/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string

		// 提前准备数据
		before func(t *testing.T)

		// 验证并且删除数据
		after   func(t *testing.T)
		reqBody string

		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				// 不需要，也就是 Redis 什么数据也没有
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
				val, err := rdb.GetDel(ctx, "phone_code:login:123xxxxxxxx").Result()
				cancel()
				assert.NoError(t, err)
				// 验证码是6位
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phone": "123xxxxxxxx",
	
}
`,

			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 0,
				Msg:  "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
				_, err := rdb.Set(ctx, "phone_code:login:123xxxxxxxx", "123456", time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
				val, err := rdb.GetDel(ctx, "phone_code:login:123xxxxxxxx").Result()
				cancel()
				assert.NoError(t, err)
				// 验证码是6位
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phone": "123xxxxxxxx",
	
}
`,

			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 0,
				Msg:  "发送成功",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			t.Log(resp)

			server.ServeHTTP(resp, req)

			var res web.Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, res)
			tc.after(t)
		})
	}
}
