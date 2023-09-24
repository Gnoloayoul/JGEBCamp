package dao

import (
	"context"
	"database/sql"
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct{
		name string

		mock func(t *testing.T) *sql.DB

		ctx context.Context
		user User

		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				// "INSERT INTO 'users' .*"
				//  这是正则表达式
				// 表示只要是 INSERT 到 users 的任意语句
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "631821745@qq.com",
					Valid: true,
				},
			},
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// "INSERT INTO 'users' .*"
				//  这是正则表达式
				// 表示只要是 INSERT 到 users 的任意语句
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
				})
				require.NoError(t, err)
				return mockDB
			},
			user: User{},
			wantErr: ErrUserDuplicate,
			},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// "INSERT INTO 'users' .*"
				//  这是正则表达式
				// 表示只要是 INSERT 到 users 的任意语句
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(errors.New("数据库错误"))
				require.NoError(t, err)
				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "631821745@qq.com",
					Valid: true,
				},
			},
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn: tc.mock(t),
				// 跳步骤
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// mock DB 不需要 ping
				DisableAutomaticPing: true,
				// 跳步骤
				SkipDefaultTransaction: true,
			})
			d := NewUserDAO(db)
			err = d.Insert(tc.ctx, tc.user)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
