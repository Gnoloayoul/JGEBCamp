package service

import (
	"context"
	"errors"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	repomocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestUserServiceIn_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct{
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入
		email string
		password string

		// 输出
		wantUser domain.User
		wantErr error
	}{
		{
			name: "登录成功",

			mock: func(ctrl *gomock.Controller) repository.UserRepository{
				usersvc := repomocks.NewMockUserRepository(ctrl)
				usersvc.EXPECT().FindByEmail(gomock.Any(), "631821745@qq.com").Return(
					domain.User{
						Email: "631821745@qq.com",
						Password: "$2a$10$exLFMrG98LXwIJ/s2/clteAH0wa5.4d2oWKlKGd3wg/plysMj8Lhm",
						Phone: "135xxxxxxxx",
						Ctime: now,
					}, nil)
				return usersvc
			},

			email: "631821745@qq.com",
			password: "hello#world12345",

			wantUser: domain.User{
				Email: "631821745@qq.com",
				Password: "$2a$10$exLFMrG98LXwIJ/s2/clteAH0wa5.4d2oWKlKGd3wg/plysMj8Lhm",
				Phone: "135xxxxxxxx",
				Ctime: now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",

			mock: func(ctrl *gomock.Controller) repository.UserRepository{
				usersvc := repomocks.NewMockUserRepository(ctrl)
				usersvc.EXPECT().FindByEmail(gomock.Any(), "631821745@qq.com").Return(
					domain.User{}, repository.ErrUserNotFound)
				return usersvc
			},

			email: "631821745@qq.com",
			password: "hello#world12345",

			wantUser: domain.User{},
			wantErr: ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",

			mock: func(ctrl *gomock.Controller) repository.UserRepository{
				usersvc := repomocks.NewMockUserRepository(ctrl)
				usersvc.EXPECT().FindByEmail(gomock.Any(), "631821745@qq.com").Return(
					domain.User{}, errors.New("someting worng"))
				return usersvc
			},

			email: "631821745@qq.com",
			password: "hello#world12345",

			wantUser: domain.User{},
			wantErr: errors.New("someting worng"),
		},
		{
			name: "密码不对",

			mock: func(ctrl *gomock.Controller) repository.UserRepository{
				usersvc := repomocks.NewMockUserRepository(ctrl)
				usersvc.EXPECT().FindByEmail(gomock.Any(), "631821745@qq.com").Return(
					domain.User{
						Email: "631821745@qq.com",
						Password: "$2a$10$exLFMrG98LXwIJ/s2/clteAH0wa5.4d2oWKlKGd3wg/plysMj8Lhm",
						Phone: "135xxxxxxxx",
						Ctime: now,
					}, nil)
				return usersvc
			},

			email: "631821745@qq.com",
			password: "121321321hello#world12345",

			wantUser: domain.User{},
			wantErr: ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestUserServiceIn_Login1(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		email string
		password string

		wantErr error
		wantUser domain.User
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository{
				usersvc := repomocks.NewMockUserRepository(ctrl)
				usersvc.EXPECT().FindByEmail(context.Background(), "631821745@qq.com").
					Return(domain.User{
						Email: "631821745@qq.com",
						Password: "$2a$10$xOagQhIGoJxrW8altcCFbudZeGGLf..rrwWokLUweoZKLLR9CC30u",
						Phone: "123XXXXXXXX",
						Ctime: now,
				}, nil)
				return usersvc
			},

			email: "631821745@qq.com",
			password: "hello#world123456",

			wantErr: nil,
			wantUser: domain.User{
				Email: "631821745@qq.com",
				Password: "$2a$10$xOagQhIGoJxrW8altcCFbudZeGGLf..rrwWokLUweoZKLLR9CC30u",
				Phone: "123XXXXXXXX",
				Ctime: now,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := NewUserService(tc.mock(ctrl))
			u, err := server.Login(context.Background(), tc.email, tc.password)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello#world123456"), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	t.Log(string(res))
}