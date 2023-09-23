package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/cache"
	cachemocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/cache/mocks"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao"
	daomocks "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx context.Context
		id  int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中, 查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存未命中：查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{
						Id: 123,
						Email: sql.NullString{
							String: "631821745@qq.com",
							Valid:  true,
						},
						Password: "xxx",
						Phone: sql.NullString{
							String: "137xxxxxxxx",
							Valid:  true,
						},
						Ctime: now.UnixMilli(),
						Utime: now.UnixMilli(),
					}, nil)

				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "631821745@qq.com",
					Phone:    "137xxxxxxxx",
					Password: "xxx",
					Ctime:    now,
				}).Return(nil)
				return d, c
			},

			ctx: context.Background(),
			id:  123,

			wantErr: nil,
			wantUser: domain.User{
				Id:       123,
				Email:    "631821745@qq.com",
				Phone:    "137xxxxxxxx",
				Password: "xxx",
				Ctime:    now,
			},
		},
		{
			name: "缓存命中, 查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存命中：查了缓存，有结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:       123,
						Email:    "631821745@qq.com",
						Phone:    "137xxxxxxxx",
						Password: "xxx",
						Ctime:    now,
					}, nil)

				d := daomocks.NewMockUserDAO(ctrl)
				return d, c
			},

			ctx: context.Background(),
			id:  123,

			wantErr: nil,
			wantUser: domain.User{
				Id:       123,
				Email:    "631821745@qq.com",
				Phone:    "137xxxxxxxx",
				Password: "xxx",
				Ctime:    now,
			},
		},
		{
			name: "缓存未命中, 查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存未命中：查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{}, errors.New("mock db 错误"))
				return d, c
			},

			ctx: context.Background(),
			id:  123,

			wantErr:  errors.New("mock db 错误"),
			wantUser: domain.User{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(tc.ctx, tc.id)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
			// 针对 go Routine 的测法
			// 不可能1秒搞不定
			time.Sleep(time.Second)
		})
	}
}

