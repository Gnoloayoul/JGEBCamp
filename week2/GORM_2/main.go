package main

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// User --> 数据表
// 定义模型
type User struct {
	gorm.Model   // 内部gorm.Model
	Name         string
	Age          sql.NullInt64 // 零值
	Birthday     *time.Time
	Email        string  `gorm:"type:varchar(100);unique_index"` // unique_index唯一记录
	Role         string  `gorm:"size:255"`                       // 设置字段大小为255
	MemberNumber *string `gorm:"unique;not null"`                // 设置会员号（member number）唯一并且不为空
	Num          int     `gorm:"AUTO_INCREMENT"`                 // 设置 num 为自增类型
	Address      string  `gorm:"index:addr"`                     // 给address字段创建名为addr的索引
	IgnoreMe     int     `gorm:"-"`                              // 忽略本字段
}

// 使用`AnimalID`作为主键
type Animal struct {
	AnimalID int64 `gorm:"primary_key"`
	Name     string
	Age      int64
}

// 将 Animal 改名
func (Animal) TableName() string {
	return "Campione"
}

func main() {
	// 连MysSQL数据库
	dsn := "acs:root278803@(119.45.240.2:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 禁用默认表名的复数形式，如果置为 true，则 `User` 的默认表名是 `user`
	// 但是现在的Gorm是没有这函数了
	//db.SingularTable(true)

	// 创建表 自动迁移 （把结构体和数据表进行对应）
	//db.AutoMigrate(&User{})
	//db.AutoMigrate(&Animal{})

}
