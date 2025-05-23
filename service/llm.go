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
	"go_logistics/model/entity"
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
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	username := c.GetString("name")
	// 第一次对话
	var history []llms.MessageContent
	isFirst := req.ChatId == ""
	if isFirst {
		history = []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, config.SystemPrompt),
		}
	} else {
		var chatService *entity.ChatService
		chatService, err = entity.GetChatByIdAndUsername(req.ChatId, username)
		if err != nil {
			config.Log.Error("获取对话信息失败", zap.Error(err))
			common.ErrorResponse(c, common.ServerError("获取对话信息失败"))
			return
		}
		history = chatService.Message
	}

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
		_, err := model.GenerateContent(c.Request.Context(), history, llms.WithStreamingFunc(streamingFunc))
		if err != nil {
			config.Log.Error("LLM 生成内容失败", zap.Error(err))
		}

		// 将 AI 回复加入历史记录
		history = append(history, llms.TextParts(llms.ChatMessageTypeAI, string(responseContent)))

		if isFirst {
			titleContent := history
			content, err = formatPrompt("帮我根据上面的内容生成一个标题，十个字以内，你只能回答我一个十个字以内标题，不需要其他内容")
			if err != nil {
				common.ErrorResponse(c, common.ServerError("生成标题失败！"))
			}
			titleContent = append(titleContent, content...)
			titleResponse, err := model.GenerateContent(c.Request.Context(), titleContent)
			if err != nil {
				common.ErrorResponse(c, common.ServerError("生成标题失败！"))
			}
			title := titleResponse.Choices[0].Content
			err = entity.InsertChat(&entity.ChatService{
				Username: username,
				Title:    title,
				Message:  history,
			})
			if err != nil {
				common.ErrorResponse(c, common.ServerError("保存聊天记录失败！"))
			}
		} else {
			err = entity.UpdateChatMessage(&entity.ChatService{
				ID:       req.ChatId,
				Username: username,
				Message:  history,
			})
			if err != nil {
				common.ErrorResponse(c, common.ServerError("更新聊天记录失败！"))
			}
		}
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

func GetChatList(c *gin.Context) {
	username := c.GetString("name")
	if username == "" {
		common.ErrorResponse(c, common.NotLogin)
		return
	}
	chats, err := entity.GetChatListByUserName(username)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponseWithData(c, chats)
}

func GetChat(c *gin.Context) {
	chatId := c.Query("chatId")
	username := c.GetString("name")
	if chatId == "" || username == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	chat, err := entity.GetChatByIdAndUsername(chatId, username)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponseWithData(c, chat)
}

func DeleteChat(c *gin.Context) {
	chatId := c.Query("chatId")
	if chatId == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	err := entity.DeleteChat(chatId)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponse(c)
}

func UpdateChatTitle(c *gin.Context) {
	title := c.Query("title")
	chatId := c.Query("chatId")
	username := c.GetString("name")
	if title == "" || chatId == "" || username == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	err := entity.UpdateChatTitle(chatId, username, title)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponse(c)
}
