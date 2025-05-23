package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"go.uber.org/zap"
	"go_logistics/common"
	"go_logistics/config"
	"go_logistics/model/dto"
	"io"
	"net/http"
)

const (
	LLMModelHunyuanLite  string = "hunyuan-lite"
	LLMModelHunyuanTurbo string = "hunyuan-turbos-latest"
	LLMModelDeepseek     string = "deepseek-reasoner"
)

func ChatLLM(c *gin.Context) {
	var req dto.LLMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		common.ErrorResponse(c, common.ParamError)
		return
	}
	content, err := formatPrompt(req.Message)
	if err != nil {
		common.ErrorResponse(c, common.ServerError("格式化提示词错误！"))
		return
	}
	// 设置流式响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	// 获取 gin.ResponseWriter 的底层 http.Flusher 接口
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		common.ErrorResponse(c, common.ServerError("流式传输不被支持！"))
		return
	}
	// 开始流式响应
	c.Stream(func(w io.Writer) bool {
		streamingFunc := func(ctx context.Context, chunk []byte) error {
			_, err := fmt.Fprintf(w, "%s", chunk)
			if err != nil {
				return err
			}

			flusher.Flush()
			return nil
		}
		var model llms.Model
		switch req.Model {
		case LLMModelHunyuanLite:
			model = config.HunyuanLite
		case LLMModelHunyuanTurbo:
			model = config.HunyuanTurbo
		case LLMModelDeepseek:
			model = config.Deepseek
		default:
			common.ErrorResponse(c, common.ServerError("不支持的模型！"))
			return false
		}
		_, err := model.GenerateContent(context.Background(), content, llms.WithStreamingFunc(streamingFunc))
		if err != nil {
			config.Log.Error("LLM 生成内容失败", zap.Error(err))
		}
		return false // 只执行一次
	})
}

// 格式化提示词
func formatPrompt(userPrompt string) (content []llms.MessageContent, err error) {

	promptTemplate := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(config.SystemPrompt, nil),
		prompts.NewHumanMessagePromptTemplate(`{{.question}}`, []string{".question"}),
	})
	prompt, err := promptTemplate.FormatMessages(map[string]any{"question": userPrompt})
	if err != nil {
		return
	}
	content = []llms.MessageContent{
		llms.TextParts(prompt[0].GetType(), prompt[0].GetContent()),
		llms.TextParts(prompt[1].GetType(), prompt[1].GetContent()),
	}
	return
}
