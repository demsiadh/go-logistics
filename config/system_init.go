package config

import (
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	MongoClient  *mongo.Client
	Log          *zap.Logger
	HunyuanLite  *openai.LLM
	HunyuanTurbo *openai.LLM
	DeepseekR1   *openai.LLM
	DeepseekV3   *openai.LLM
	SystemPrompt string
)

func init() {
	initEnvConfig()
	initLog()
	initMongoDB()
	initLLM()
	initPinecone()
}
