package main

import (
	"bytes"
	"github.com/spf13/viper"
)

// initViperReader
// 直接从内存里取配置
// 千万别漏了这个 viper.SetConfigType ，写在内存里的内容总得有个格式
func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dns: "XXXXX"
`

	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

func main() {
	initViperReader()

	type Config struct {
		DNS string `yaml:dns`
	}
	var cfg Config
	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}
	println(cfg.DNS)
}
