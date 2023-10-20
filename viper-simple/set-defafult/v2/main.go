package main

import (
	"github.com/spf13/viper"
)

func initViper() {
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
		Name string `yaml:name`
	}
	cfg := Config{
		Name: "x123456789",
	}
	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}

	println(cfg.Name)

}
