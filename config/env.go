package config

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	MongodbUri        string
	MongodbDatabase   string
	MongodbUsername   string
	MongodbPassword   string
	MongodbAuthSource string
	SecretKey         string
	HunyuanApiKey     string
	DeepseekApiKey    string
	AliBaiLianApiKey  string
	PineconeApiKey    string
)

func initEnvConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	MongodbUri = os.Getenv("MONGODB_URI")
	MongodbDatabase = os.Getenv("MONGODB_DATABASE")
	MongodbUsername = os.Getenv("MONGODB_USERNAME")
	MongodbPassword = os.Getenv("MONGODB_PASSWORD")
	MongodbAuthSource = os.Getenv("MONGODB_AUTH_SOURCE")
	SecretKey = os.Getenv("SECRET_KEY")
	HunyuanApiKey = os.Getenv("HUNYUAN_API_KEY")
	DeepseekApiKey = os.Getenv("DEEPSEEK_API_KEY")
	AliBaiLianApiKey = os.Getenv("ALI_BAILIAN_API_KEY")
	PineconeApiKey = os.Getenv("PINECONE_API_KEY")
}
