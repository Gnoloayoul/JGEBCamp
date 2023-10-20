package main

import (
	"github.com/spf13/viper"
)

func initViperV2() {
	// 指定配置文件路径
	viper.SetConfigFile("config/dev.yaml")
	// 读配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	initViperV2()
	msg := viper.GetString("db.mysql.dsn")
	println(msg)
}

