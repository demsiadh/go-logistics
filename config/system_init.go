package config

import (
	"go.uber.org/zap"
)

var (
	initStep int
)

func init() {
	initLog()
	initEnvConfig()
	initMongoDB()
	initLLM()
	initPinecone()
}

// 处理初始化成功日志
func handleSuccess(msg string) {
	initStep++
	Log.Info(msg, zap.Int("initStep", initStep))
}

// 处理初始化失败日志
func handleError(msg string, err error) {
	if err != nil {
		initStep++
		Log.Error(msg, zap.Int("initStep", initStep), zap.Error(err))
		panic(err)
	}
}
