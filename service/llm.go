package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	easyllm "github.com/soryetong/go-easy-llm"
	"github.com/soryetong/go-easy-llm/easyai/chatmodule"
	"go_logistics/common"
	"go_logistics/config"
	"net/http"
)

var LLMConfig = easyllm.DefaultConfigWithSecret("666", "666", chatmodule.ChatTypeHunYuan)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

type UserInput struct {
	Message string `json:"message"`
}

func ChatLLM(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	defer conn.Close()

	for {
		client := easyllm.NewChatClient(LLMConfig)
		if err := handleWebSocket(conn, client); err != nil {
			config.Log.Info(err.Error())
			break
		}
	}
}

func handleWebSocket(conn *websocket.Conn, client *easyllm.ChatClient) error {
	var input UserInput
	_, reader, err := conn.NextReader()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			config.Log.Info("Client closed the connection gracefully")
			return err
		}
		config.Log.Info("Failed to read user message")
		return err
	}

	if err := json.NewDecoder(reader).Decode(&input); err != nil {
		config.Log.Info("Failed to decode user message:")
		return err
	}

	presetPrompt := "你是一个物流管理系统的客服助手，可以处理以下功能：\n1. 订单管理：查询订单状态、跟踪订单状态。\n2. 用户管理：管理用户账号。\n3. 车辆管理：查询车辆信息。\n4. 线路管理：管理运输线路。\n5. 营业网点管理：查询网点信息。\n请根据用户的问题提供专业解答。"
	fullMessage := fmt.Sprintf("%s\n%s", presetPrompt, input.Message)

	resp, err := client.StreamChat(context.Background(), &chatmodule.ChatRequest{
		Model:   "hunyuan-lite",
		Message: fullMessage,
	})
	if err != nil {
		config.Log.Info("Model call failed:")
		return err
	}

	for content := range resp {
		msg := WSMessage{Content: content.Content, Done: false}
		if err := conn.WriteJSON(msg); err != nil {
			config.Log.Info("WebSocket write failed:")
			return err
		}
	}

	// Send completion signal
	conn.WriteJSON(WSMessage{Done: true})
	return nil
}
