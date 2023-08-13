package domain

import (
	"time"
)

// User
// (领域对象)
// BO(business object)
type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}

// UserInfo
// 用户信息
type UserInfo struct {
	Id int64
	NickName string
	Birthday string
	Profile string
}
