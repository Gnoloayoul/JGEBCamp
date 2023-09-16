package web

import (
	"bytes"
	"context"
	"errors"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	svcmocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/service/mocks"
	ijwt "github.com/Gnoloayoul/JGEBCamp/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
	} {
		{},
	}

	// 构建请求
	// 用http.NewRequest
	req, err := http.NewRequest(http.MethodPost, "users/signup", bytes.NewBuffer([]byte(`
{
	"email": "123@qq.com",
	"password": "123456"
}
`)))
	// 这里假定url绝对不会错，不会引发err
	require.NoError(t, err)
	// 可以继续使用 req

	// 构建响应
	// 用http.ResponseWriter
	// httptest.NewRecorder
	//func NewRecorder() *ResponseRecorder {
	//	return &ResponseRecorder{
	//	HeaderMap: make(http.Header),   http头
	//	Body:      new(bytes.Buffer),   【验证】写进去的数据
	//	Code:      200,                 【验证】http响应码
	//}
	//}
	resp := httptest.NewRecorder()

	// Handler要这样初始化？
	// h := NewUserHandler(nil, nil)
	// 传nil进去，不就直接爆err吗？


	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){

		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersvc := svcmocks.NewMockUserService(ctrl)

	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).return(errors.New("mock error"))

	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})
	t.Log(err)
}