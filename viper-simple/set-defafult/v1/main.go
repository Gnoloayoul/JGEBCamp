package main

import (
	"github.com/spf13/viper"
)

func initViper() {
	// 读为默认值
	viper.SetDefault("db.mysql.name", "133")

	// 指定配置文件路径
	viper.SetConfigFile("config/dev.yaml")
	// 读配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	initViper()

	type Config struct {
		DSN string `yaml:dsn`
		Name string `yaml:name`
	}
	var cfg Config

	type Config2 struct {
		Name string `yaml:name`
	}
	var cfg2 Config2

	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}

	err = viper.UnmarshalKey("db.mysql", &cfg2)
	if err != nil {
		panic(err)
	}
	println(cfg.DSN, " ", cfg.Name)
	println(cfg2.Name)
}
