package config

import (
	"github.com/tmc/langchaingo/llms/openai"
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

var (
	HunyuanLite  *openai.LLM
	HunyuanTurbo *openai.LLM
	DeepseekR1   *openai.LLM
	DeepseekV3   *openai.LLM
	SystemPrompt string
)

func initLLM() {
	initHunyuanLite()
	initHunyuanTurbo()
	initDeepseekR1()
	initDeepseekV3()
	initSystemPrompt()
}

// 初始化hunyuanLite
func initHunyuanLite() {
	var err error
	// 初始化hunyuanLite
	HunyuanLite, err = openai.New(
		openai.WithBaseURL(HunyuanBaseUrl),
		openai.WithModel(HunyuanLiteModel),
		openai.WithToken(HunyuanApiKey),
	)
	handleError("初始化hunyuanLite失败！", err)
	handleSuccess("初始化hunyuanLite成功！")
}

// 初始化hunyuanTurbo
func initHunyuanTurbo() {
	var err error
	HunyuanTurbo, err = openai.New(
		openai.WithBaseURL(HunyuanBaseUrl),
		openai.WithModel(HunyuanTurboModel),
		openai.WithToken(HunyuanApiKey),
	)
	handleError("初始化hunyuanTurbo失败！", err)
	handleSuccess("初始化hunyuanTurbo成功！")
}

// 初始化deepseekr1
func initDeepseekR1() {
	var err error
	DeepseekR1, err = openai.New(
		openai.WithBaseURL(DeepseekBaseUrl),
		openai.WithModel(DeepseekR1Model),
		openai.WithToken(DeepseekApiKey),
	)
	handleError("初始化deepseekR1失败！", err)
	handleSuccess("初始化deepseekR1成功！")
}

// 初始化deepseekv3
func initDeepseekV3() {
	var err error
	DeepseekV3, err = openai.New(
		openai.WithBaseURL(DeepseekBaseUrl),
		openai.WithModel(DeepseekV3Model),
		openai.WithToken(DeepseekApiKey),
	)
	handleError("初始化deepseekV3失败！", err)
	handleSuccess("初始化deepseekV3成功！")
}

func initSystemPrompt() {
	// 初始化系统提示词
	pwd, err := os.Getwd()
	handleError("获取当前工作目录失败！", err)

	systemPromptPath := filepath.Join(pwd, "static", "system_prompt.txt")
	SystemPromptByte, err := os.ReadFile(systemPromptPath)
	handleError("读取系统提示词失败！", err)

	SystemPrompt = string(SystemPromptByte)
	handleSuccess("初始化系统提示词成功！")
}
