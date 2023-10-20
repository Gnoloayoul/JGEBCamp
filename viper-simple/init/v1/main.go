package main

import (
	"github.com/spf13/viper"
)

func initViperv1() {

	// 指定配置文件名 dev
	viper.SetConfigName("dev")
	// 指定配置文件格式 yaml
	viper.SetConfigType("yaml")
	// 指定配置文件名的查找路径 /config
	// 可以添加多个路径，一行一个
	viper.AddConfigPath("./config")
	// 读这配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	initViperv1()

	// 用法1
	msg := viper.GetString("db.mysql.dsn")
	println(msg)

	// 用法2 大明推荐
	type Config struct {
		// 在这里直接指定了在yaml文件的 dsn 字段
		DSN string `yaml:dsn`
	}
	var cfg Config
	// viper.UnmarshalKey
	// 按 key 解析提取所需
	// key 是目标字段 dsn 的上一级（从哪里开始找）
	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}
	println("way2: ", cfg.DSN)
}

