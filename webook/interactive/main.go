package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)


func main() {
	initViperV1()
	app := InitAPP()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	err := app.server.Serve()
	log.Println(err)
}

func initViperV1() {
	cfile := pflag.String("config",
		"config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()

	// 只告诉你文件变了
	// 比较好的设计： 能告诉你变之前和变之后的数据
	// 更好的设计：直接告诉你差异在哪里
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
