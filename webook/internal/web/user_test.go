package web

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	svcmocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/service/mocks"
	ijwt "github.com/Gnoloayoul/JGEBCamp/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandlerV1_LoginSMSV1(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.UserService,
			service.CodeService, ijwt.Handler)
		reqBuilder func(t *testing.T) *http.Request
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService,
				service.CodeService, ijwt.Handler) {
				usersvc := svcmocks.NewMockCodeService(ctrl)
				usersvc.EXPECT().LoginSMSV1(gomock.Any()).
					Return(nil)
				codesvc := svcmocks.


			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usersvc, codesvc, jwthdl := tc.mock(ctrl)
			hdl := NewUserHandlerV1(usersvc, codesvc, jwthdl)

			// 准备gin引擎
			server := gin.Default()
			// 注册路由（模拟）
			hdl.RegisterRoutesV1(server)
			// 准备请求
			req := tc.reqBuilder(t)
			// 准备记录响应
			recorder := httptest.NewRecorder()
			// 执行
			server.ServeHTTP(recorder, req)
			// 断言
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body)
		})
	}
}
