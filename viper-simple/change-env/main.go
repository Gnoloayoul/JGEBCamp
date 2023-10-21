package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func initViper() {
	cfile := pflag.String("config", "config/config.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	initViper()

	type Message struct {
		Name string `yaml:name`
	}
	var msg Message
	err := viper.UnmarshalKey("Table", &msg)
	if err != nil {
		panic(err)
	}
	println(msg.Name)
}
