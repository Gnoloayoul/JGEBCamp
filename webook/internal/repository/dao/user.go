package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
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

func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	//err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("'id' = ?", id).First(&u).Error
	return u, err
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime, u.Ctime = now, now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突 or 手机号码冲突
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *UserDAO) Edit(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now

	return dao.db.WithContext(ctx).Model(&u).Updates(User{
		Utime:    u.Utime,
		NickName: u.NickName,
		Birthday: u.Birthday,
		Info:     u.Info,
	}).Error
}

//func FindRecordByID(ctx context.Context, db *gorm.DB, userID string) (Record, error) {
//	var record Record
//	err := db.WithContext(ctx).Where("user_id = ?", userID).First(&record).Error
//	if err != nil {
//		return record, err
//	}
//	return record, nil
//}

func (dao *UserDAO) Profile(ctx context.Context, u User) (User, error) {
	key := strconv.FormatInt(u.Id, 10)
	err := dao.db.WithContext(ctx).Where("id = ?", key).First(&u).Error
	// 测试打印 取u之前
	//fmt.Printf("\nform dao--u: %#v", u)
	resUser := User{
		Id:       u.Id,
		Email:    u.Email,
		NickName: u.NickName,
		Birthday: u.Birthday,
		Info:     u.Info,
	}
	// 测试打印 取u之前
	//fmt.Printf("\nform dao--resUser: %#v", resUser)
	return resUser, err
}

// User
// 直接对应数据库表结构
type User struct {
	Id int64 `gorm:"primaryKer,autoIncrement"`
	// 用户唯一标识
	Email    sql.NullString `gorm:"unique"`
	Password string

	// 唯一索引允许有多个空值
	// 但是不能有多个 ""
	Phone sql.NullString `gorm:"unique"`
	// 最大问题就是，你要解引用
	// 你要判空
	//Phone *string

	// 往这面加

	// 创建时间(毫秒)
	Ctime int64
	// 更新时间(毫秒)
	Utime int64

	NickName string
	Birthday string
	Info     string
}
