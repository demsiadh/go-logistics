package dto

type LLMRequest struct {
	Model   string `json:"model"`
	Message string `json:"message"`
	ChatId  string `json:"chatId"`
	IsRAG   bool   `json:"isRAG"`
	IsAgent bool   `json:"isAgent"`
}
