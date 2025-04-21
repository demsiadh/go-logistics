package service

import (
	"context"
	"encoding/json"
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
	Model   string `json:"model"`
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

	resp, err := client.StreamChat(context.Background(), &chatmodule.ChatRequest{
		Model:   input.Model,
		Message: input.Message,
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
