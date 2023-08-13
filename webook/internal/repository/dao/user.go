package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime, u.Ctime = now, now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

// User
// 直接对应数据库表结构
type User struct {
	Id int64 `gorm:"primaryKer,autoIncrement"`
	// 用户唯一标识
	Email    string `gorm:"unique"`
	Password string

	// 创建时间(毫秒)
	Ctime int64
	// 更新时间(毫秒)
	Utime int64
}

// UserInfo
// 用户信息
type UserInfo struct {
	Id int64 `gorm:"primaryKer,autoIncrement,unique"`
	NickName string
	Birthday string
	Profile string
}