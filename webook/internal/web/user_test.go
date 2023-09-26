package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	_ "github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	svcmocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.UserService

		reqBody string

		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "631821745@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				// 注册成功是 return nil
				return usersvc
			},
			reqBody: `
{
	"email": "631821745@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对， bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				// 注册成功是 return nil
				return usersvc
			},
			reqBody: `
{
	"email": "631821745@qq.com",
	"password": "hello#world123",
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "631821745@qq",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "输入的邮箱格式不对",
		},
		{
			name: "两次输入密码不匹配",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "631821745@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello111#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)

				// 注册成功是 return nil
				return usersvc
			},
			reqBody: `
{
	"email": "631821745@qq.com",
	"password": "world123",
	"confirmPassword": "world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，包含数字、特殊字符、字母",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "745@qq.com",
					Password: "hello#world123",
				}).Return(service.ErrUserDuplicateEmail)
				// 注册成功是 return nil
				return usersvc
			},
			reqBody: `
{
	"email": "745@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "631821745@qq.com",
					Password: "hello#world123",
				}).Return(errors.New("随便一个error"))
				// 注册成功是 return nil
				return usersvc
			},
			reqBody: `
{
	"email": "631821745@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()

			// 没用上 CodeSvc
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			// 构建请求
			// 用http.NewRequest
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			// 这里假定url绝对不会错，不会引发err
			require.NoError(t, err)
			// 可以继续使用 req
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")

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
			t.Log(resp)

			// HTTP 请求进入 GIN 框架的入口
			// 当这样调的时候， GIN 就是会处理这个结果
			// 响应也会写回 resp 里
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersvc := svcmocks.NewMockUserService(ctrl)

	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))

	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})
	t.Log(err)
}

func TestUserHandler_SignUp1(t *testing.T) {
	// TODO: 用户登录时的邮箱解析失败
	// TODO: 用户登录时的密码解析失败
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "631821745@qq.com",
					Password: "a123454@123214",
				}).Return(nil)
				return usersvc
			},
			reqBody: `
{
	"email": "631821745@qq.com",
	"password": "a123454@123214",
	"confirmPassword": "a123454@123214"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "631821745@qq.com",
					Password: "a123454@123214",
				}).Return(errors.New("随便一个error"))
				return usersvc
			},
			reqBody: `{
				"email": "631821745@qq.com",
				"password": "a123454@123214",
				"confirmPassword": "a123454@123214"
}`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
		{
			name: "输入的邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `{
				"email": "6@qq",
				"password": "a123454@123214",
				"confirmPassword": "a123454@123214"
}`,
			wantCode: http.StatusOK,
			wantBody: "输入的邮箱格式不对",
		},
		{
			name: "两次输入的密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `{
				"email": "631821745@qq.com",
				"password": "a123454@123214",
				"confirmPassword": "a12@123214"
}`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码必须大于8位，包含数字、特殊字符、字母",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `{
				"email": "631821745@qq.com",
				"password": "1@123214",
				"confirmPassword": "1@123214"
}`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，包含数字、特殊字符、字母",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "631821745@qq.com",
					Password: "a123454@123214",
				}).Return(service.ErrUserDuplicateEmail)
				return usersvc
			},
			reqBody: `{
				"email": "631821745@qq.com",
				"password": "a123454@123214",
				"confirmPassword": "a123454@123214"
}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()

			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			t.Log(resp)

			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestUserHandler_LoginSMS(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBody string

		wantCode int
		wantBody Result
	}{
		{
			name: "短信登录成功",
			mock: func(ctrl *gomock.Controller) (service.UserService,service.CodeService) {
				us := svcmocks.NewMockUserService(ctrl)
				cs := svcmocks.NewMockCodeService(ctrl)

				cs.EXPECT().Verify(gomock.Any(), "login", "137XXXXXXXX", "1234567").
					Return(true, nil)
				us.EXPECT().FindOrCreate(gomock.Any(), "137XXXXXXXX").
					Return(domain.User{
						Phone: "137XXXXXXXX",
				}, nil)

				return us, cs
			},
			reqBody:
`
{
	"phone": "137XXXXXXXX",
	"code": "1234567"
}
`,


			wantCode: http.StatusOK,
			wantBody: Result{
				Msg: "验证码校验通过",
			},
		},
		{
			name: "参数错误，没能 bind 上",
			mock: func(ctrl *gomock.Controller) (service.UserService,service.CodeService) {
				us := svcmocks.NewMockUserService(ctrl)
				cs := svcmocks.NewMockCodeService(ctrl)

				return us, cs
			},
			reqBody: "",

			wantCode: http.StatusBadRequest,
			wantBody: Result{},
		},
		{
			name: "Verift 系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService,service.CodeService) {
				us := svcmocks.NewMockUserService(ctrl)
				cs := svcmocks.NewMockCodeService(ctrl)

				cs.EXPECT().Verify(gomock.Any(), "login", "137XXXXXXXX", "1234567").
					Return(true, errors.New("Verift 系统错误"))

				return us, cs
			},
			reqBody:
			`
{
	"phone": "137XXXXXXXX",
	"code": "1234567"
}
`,


			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "Verify 验证码有误",
			mock: func(ctrl *gomock.Controller) (service.UserService,service.CodeService) {
				us := svcmocks.NewMockUserService(ctrl)
				cs := svcmocks.NewMockCodeService(ctrl)

				cs.EXPECT().Verify(gomock.Any(), "login", "137XXXXXXXX", "1234567").
					Return(false, nil)

				return us, cs
			},
			reqBody:
			`
{
	"phone": "137XXXXXXXX",
	"code": "1234567"
}
`,


			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 4,
				Msg:  "验证码有误",
			},
		},
		{
			name: "FindORCreate 失败",
			mock: func(ctrl *gomock.Controller) (service.UserService,service.CodeService) {
				us := svcmocks.NewMockUserService(ctrl)
				cs := svcmocks.NewMockCodeService(ctrl)

				cs.EXPECT().Verify(gomock.Any(), "login", "137XXXXXXXX", "1234567").
					Return(true, nil)
				us.EXPECT().FindOrCreate(gomock.Any(), "137XXXXXXXX").
					Return(domain.User{
						Phone: "137XXXXXXXX",
					}, errors.New("mock FindOrCreate error"))

				return us, cs
			},
			reqBody:
			`
{
	"phone": "137XXXXXXXX",
	"code": "1234567"
}
`,


			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()

			us, cs := tc.mock(ctrl)
			uh := NewUserHandler(us, cs)
			uh.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/loginsms", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			t.Log(resp)

			server.ServeHTTP(resp, req)

			// 从响应中取出 Result ，需要使用 json 解码
			var resultData Result
			_ = json.NewDecoder(resp.Body).Decode(&resultData)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resultData)
		})
	}
}

func TestUserHandler_LoginSMS2(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBody string

		wantCode int
		wantBody Result
	}{
		{
			name: "短信设置 JWT 失败",
			mock: func(ctrl *gomock.Controller) (service.UserService,service.CodeService) {
				us := svcmocks.NewMockUserService(ctrl)
				cs := svcmocks.NewMockCodeService(ctrl)

				cs.EXPECT().Verify(gomock.Any(), "login", "137XXXXXXXX", "1234567").
					Return(true, nil)
				us.EXPECT().FindOrCreate(gomock.Any(), "137XXXXXXXX").
					Return(domain.User{
						Id: 1,
						Phone: "137XXXXXXXX",
					}, nil)
				return us, cs
			},
			reqBody:
			`
{
	"phone": "137XXXXXXXX",
	"code": "1234567"
}
`,

			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()

			us, cs := tc.mock(ctrl)
			uh := NewUserHandler(us, cs)
			uh.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/loginsms", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			t.Log(resp)

			server.ServeHTTP(resp, req)

			// 从响应中取出 Result ，需要使用 json 解码
			var resultData Result
			_ = json.NewDecoder(resp.Body).Decode(&resultData)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resultData)
		})
	}
}