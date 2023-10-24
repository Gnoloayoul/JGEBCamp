package main

import (
	"errors"
	"go.uber.org/zap"
)

func initLogger() {
	// 新建 logger 对象
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// 将 logger 放入全局，使之生效
	// 这句不能没有
	zap.ReplaceGlobals(logger)

	zap.L().Info("before: Hello")
	zap.L().Info("after: Hello")

	type Ex struct {
		Name string `json:"name"`
	}
	zap.L().Info("This is ex-code",
		zap.Error(errors.New("This is an error")),
		zap.Int64("id", 123),
		zap.Any("one struct", Ex{Name: "Herry"}))

	// 敏感数据不能进入日志
	// 万不得已要这样用
	// 也要先把敏感数据给加密了
	// 又或者将其脱敏， 例如 123******456
	zap.L().Debug("手机号", zap.Error(err))
}

func main() {
	initLogger()
}
