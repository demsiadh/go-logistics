package config

import (
	"github.com/tmc/langchaingo/llms/openai"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func initLLM() {
	var err error
	// 初始化hunyuanLite
	HunyuanLite, err = openai.New(
		openai.WithBaseURL("https://api.hunyuan.cloud.tencent.com/v1"),
		openai.WithModel("hunyuan-lite"),
		openai.WithToken(os.Getenv("HUNYUAN_API_KEY")),
	)
	if err != nil {
		Log.Error("初始化hunyuanLite失败！", zap.Error(err))
		panic(err)
	}
	// 初始化hunyuanTurbo
	HunyuanTurbo, err = openai.New(
		openai.WithBaseURL("https://api.hunyuan.cloud.tencent.com/v1"),
		openai.WithModel("hunyuan-turbos-latest"),
		openai.WithToken(os.Getenv("HUNYUAN_API_KEY")),
	)
	if err != nil {
		Log.Error("初始化hunyuanTurbo失败！", zap.Error(err))
		panic(err)
	}
	// 初始化deepseek
	DeepseekR1, err = openai.New(
		openai.WithBaseURL("https://api.deepseek.com"),
		openai.WithModel("deepseek-reasoner"),
		openai.WithToken(os.Getenv("DEEPSEEK_API_KEY")),
	)
	if err != nil {
		Log.Error("初始化DeepSeekR1失败！", zap.Error(err))
		panic(err)
	}
	DeepseekV3, err = openai.New(
		openai.WithBaseURL("https://api.deepseek.com"),
		openai.WithModel("deepseek-chat"),
		openai.WithToken(os.Getenv("DEEPSEEK_API_KEY")),
	)
	if err != nil {
		Log.Error("初始化DeepSeekV3失败！", zap.Error(err))
		panic(err)
	}
	// 初始化系统提示词
	pwd, err := os.Getwd()
	if err != nil {
		Log.Error("获取当前工作目录失败！", zap.Error(err))
		panic(err)
	}
	systemPromptPath := filepath.Join(pwd, "static", "system_prompt.txt")
	SystemPromptByte, err := os.ReadFile(systemPromptPath)
	if err != nil {
		Log.Error("读取系统提示语失败！", zap.Error(err))
		panic(err)
	}
	SystemPrompt = string(SystemPromptByte)
	Log.Info("初始化LLM成功！")
}
