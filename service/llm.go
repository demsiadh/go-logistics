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
	"sync"
)

const (
	LLMModelHunyuanLite  string = "hunyuan-lite"
	LLMModelHunyuanTurbo string = "hunyuan-turbos-latest"
	LLMModelDeepseek     string = "deepseek-reasoner"
)

// 定义互斥锁
var mu sync.Mutex
var chatHistoryMap = make(map[string][]llms.MessageContent) // key: 用户ID 或 会话ID

func ChatLLM(c *gin.Context) {
	var req dto.LLMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		common.ErrorResponse(c, common.ParamError)
		return
	}
	username := c.GetString("name")
	// 获取或初始化该用户的对话历史
	mu.Lock()
	history, exists := chatHistoryMap[username]
	if !exists {
		history = []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, config.SystemPrompt),
		}
	}
	mu.Unlock()

	// 格式化当前用户输入
	content, err := formatPrompt(req.Message)
	if err != nil {
		common.ErrorResponse(c, common.ServerError("格式化提示词错误！"))
		return
	}

	// 将用户输入追加到历史记录中
	history = append(history, content...)

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
		var responseContent []byte

		streamingFunc := func(ctx context.Context, chunk []byte) error {
			_, err := fmt.Fprintf(w, "%s", chunk)
			if err != nil {
				return err
			}
			responseContent = append(responseContent, chunk...)
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

		// 使用完整的对话历史调用模型
		_, err := model.GenerateContent(context.Background(), history, llms.WithStreamingFunc(streamingFunc))
		if err != nil {
			config.Log.Error("LLM 生成内容失败", zap.Error(err))
		}

		// 将 AI 回复加入历史记录
		history = append(history, llms.TextParts(llms.ChatMessageTypeAI, string(responseContent)))

		// 更新历史记录
		mu.Lock()
		chatHistoryMap[username] = history
		mu.Unlock()

		return false // 只执行一次
	})
}

// 格式化提示词
func formatPrompt(userPrompt string) (content []llms.MessageContent, err error) {
	promptTemplate := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewHumanMessagePromptTemplate(`{{.question}}`, []string{".question"}),
	})
	prompt, err := promptTemplate.FormatMessages(map[string]any{"question": userPrompt})
	if err != nil {
		return
	}
	content = []llms.MessageContent{
		llms.TextParts(prompt[0].GetType(), prompt[0].GetContent()),
	}
	return
}
