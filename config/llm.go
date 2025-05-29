package config

import (
	"github.com/tmc/langchaingo/llms/openai"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

const (
	HunyuanBaseUrl    = "https://api.hunyuan.cloud.tencent.com/v1"
	HunyuanLiteModel  = "hunyuan-lite"
	HunyuanTurboModel = "hunyuan-turbos-latest"
	DeepseekBaseUrl   = "https://api.deepseek.com"
	DeepseekV3Model   = "deepseek-chat"
	DeepseekR1Model   = "deepseek-reasoner"
)

func initLLM() {
	var err error
	// 初始化hunyuanLite
	HunyuanLite, err = openai.New(
		openai.WithBaseURL(HunyuanBaseUrl),
		openai.WithModel(HunyuanLiteModel),
		openai.WithToken(HunyuanApiKey),
	)
	if err != nil {
		Log.Error("初始化hunyuanLite失败！", zap.Error(err))
		panic(err)
	}
	// 初始化hunyuanTurbo
	HunyuanTurbo, err = openai.New(
		openai.WithBaseURL(HunyuanBaseUrl),
		openai.WithModel(HunyuanTurboModel),
		openai.WithToken(HunyuanApiKey),
	)
	if err != nil {
		Log.Error("初始化hunyuanTurbo失败！", zap.Error(err))
		panic(err)
	}
	// 初始化deepseek
	DeepseekR1, err = openai.New(
		openai.WithBaseURL(DeepseekBaseUrl),
		openai.WithModel(DeepseekR1Model),
		openai.WithToken(DeepseekApiKey),
	)
	if err != nil {
		Log.Error("初始化DeepSeekR1失败！", zap.Error(err))
		panic(err)
	}
	DeepseekV3, err = openai.New(
		openai.WithBaseURL(DeepseekBaseUrl),
		openai.WithModel(DeepseekV3Model),
		openai.WithToken(DeepseekApiKey),
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
