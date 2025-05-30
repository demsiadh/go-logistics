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
	"time"
)

const (
	MaxContextMessage int = 20
	MaxHistoryMessage int = 100
)

func ChatLLM(c *gin.Context) {
	var req dto.LLMRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	username := c.GetString("name")
	chatService, err := entity.GetChatByIdAndUsername(req.ChatId, username)
	if err != nil {
		common.ErrorResponse(c, common.ServerError("获取对话信息失败！"))
		return
	}
	isFirst := len(chatService.Message) == 0

	var totalHistory []llms.MessageContent
	var contextHistory []llms.MessageContent
	if isFirst {
		totalHistory = []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, config.SystemPrompt),
		}
		contextHistory = totalHistory
	} else {
		totalHistory = chatService.Message
		if len(totalHistory) < MaxContextMessage {
			contextHistory = totalHistory
		} else {
			contextHistory = totalHistory[len(totalHistory)-(MaxContextMessage-2):]
		}
	}

	// 格式化当前用户输入
	content, err := formatPrompt(req.Message, req.IsRAG)
	if err != nil {
		config.Log.Error("111", zap.Error(err))
		common.ErrorResponse(c, common.ServerError("格式化提示词错误！"))
		return
	}

	// 将用户输入追加到历史记录中
	contextHistory = append(contextHistory, content...)
	totalHistory = append(totalHistory, content...)

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
		case config.HunyuanLiteModel:
			model = config.HunyuanLite
		case config.HunyuanTurboModel:
			model = config.HunyuanTurbo
		case config.DeepseekR1Model:
			model = config.DeepseekR1
		case config.DeepseekV3Model:
			model = config.DeepseekV3
		default:
			common.ErrorResponse(c, common.ServerError("不支持的模型！"))
			return false
		}

		// 使用完整的对话历史调用模型
		_, err := model.GenerateContent(c.Request.Context(), contextHistory, llms.WithStreamingFunc(streamingFunc))
		if err != nil {
			config.Log.Error("LLM 生成内容失败", zap.Error(err))
		}

		// 将 AI 回复加入历史记录
		totalHistory = append(totalHistory, llms.TextParts(llms.ChatMessageTypeAI, string(responseContent)))

		if isFirst {
			titleContent := totalHistory[1:]
			content, err = formatPrompt("帮我根据上面的内容生成一个标题，十个字以内，你只能回答我一个十个字以内标题，不需要其他内容，如果无法生成标题，就输出无标题", false)
			if err != nil {
				common.ErrorResponse(c, common.ServerError("生成标题失败！"))
			}
			titleContent = append(titleContent, content...)
			titleResponse, err := model.GenerateContent(c.Request.Context(), titleContent)
			if err != nil {
				common.ErrorResponse(c, common.ServerError("生成标题失败！"))
			}
			title := titleResponse.Choices[0].Content
			if len(title) > 10 {
				title = "无标题"
			}
			err = entity.UpdateChat(&entity.ChatService{
				ID:       req.ChatId,
				Username: username,
				Title:    title,
				Message:  totalHistory,
			})
			if err != nil {
				common.ErrorResponse(c, common.ServerError("保存聊天记录失败！"))
			}
		} else {
			if len(totalHistory) > MaxHistoryMessage {
				totalHistory = totalHistory[len(totalHistory)-MaxHistoryMessage:]
			}
			err = entity.UpdateChat(&entity.ChatService{
				ID:       req.ChatId,
				Username: username,
				Message:  totalHistory,
			})
			if err != nil {
				common.ErrorResponse(c, common.ServerError("更新聊天记录失败！"))
			}
		}
		return false // 只执行一次
	})
}

func formatPrompt(userPrompt string, isRAG bool) (content []llms.MessageContent, err error) {
	var contextText string
	if isRAG {
		textSegments, err := entity.SearchAndExtractContents(context.Background(), userPrompt, 1)
		if err != nil {
			return nil, err
		}
		contextText = ""
		for _, seg := range textSegments {
			contextText += seg + "\n"
		}
	}

	promptTemplate := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewGenericMessagePromptTemplate("generic", "以下是可能相关的背景信息：\n{{.context}}", []string{"context"}),
		prompts.NewHumanMessagePromptTemplate("{{.input}}", []string{"input"}),
	})

	messages, err := promptTemplate.FormatMessages(map[string]interface{}{
		"context": contextText,
		"input":   userPrompt,
	})
	if err != nil {
		return nil, err
	}

	content = make([]llms.MessageContent, len(messages))
	for i, msg := range messages {
		content[i] = llms.TextParts(msg.GetType(), msg.GetContent())
	}
	return content, nil
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
	err := entity.UpdateChat(&entity.ChatService{
		ID:       chatId,
		Title:    title,
		Username: username,
	})
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponse(c)
}

func CreateChat(c *gin.Context) {
	username := c.GetString("name")
	title := "新对话" + time.Now().Format("2006-01-02 15:04:05")
	chatId, err := entity.InsertChat(&entity.ChatService{
		Title:    title,
		Username: username,
	})
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponseWithData(c, chatId)
}
