package service

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUserServiceIn_Login(t *testing.T) {
	testCases := []struct{
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入
		ctx context.Context
		email string
		password string

		wantUser domain.User
		wantErr error
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewUserService(tc.mock(ctrl))
			svc.Login()
		})
	}
}
