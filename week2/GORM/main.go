// GORM
// Go + ORM
// ORM: Object对象（像结构体实例） + Relational关系（关系数据库） + Mapping映射
// 数据表 -- 结构体
// 数据行 -- 结构体实例
// 字段  -- 结构体字段

// 准备mysql环境
// $ sudo docker run --name mysql8029 -p 13306:3306 -e MYSQL_ROOT_PASSWORD=root1234 -d mysql:8.0.29
// 以mysql 8.0.29版本
// 以用户名 mysql8029 ， 密码为 root1234 ， 以本地端口 13306 ,跑起来

// 登陆mysql
// $ sudo docker run -it --network host --rm mysql mysql -h127.0.0.1 -P13306 --default-character-set=utf8mb4 -uroot -p
// 会要求输入密码， 密码就是上面设定好的 root1234

// CREATE DATABASE db1; 手动创建名为 db1 的数据库
// 可用 show databases; 查看当前数据库

// 配备能远程访问的SQL用户
// CREATE USER 'username'@'%' IDENTIFIED BY 'password';
// GRANT ALL PRIVILEGES ON *.* TO 'username'@'%';
// FLUSH PRIVILEGES;
// 替换掉上面的 username 与 password

package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// UserInfo --> 数据表
type UserInfo struct {
	Id     uint
	Name   string
	Gender string
	Hobby  string
}

func main() {
	// 连MysSQL数据库
	dsn := "acs:root278803@(119.45.240.2:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 创建表 自动迁移 （把结构体和数据表进行对应）
	db.AutoMigrate(&UserInfo{})

	// 创建数据行
	u1 := UserInfo{1, "qimi", "男", "足球"}
	db.Create(&u1)

	time.Sleep(time.Second)
}
