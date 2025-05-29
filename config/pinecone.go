package config

import (
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores/pinecone"
)

const (
	EmbeddingModel    = "text-embedding-v3"
	EmbeddingBaseUrl  = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	PineconeNamespace = "langchain4j-demo-index"
	PineconeHost      = "https://langchain4j-demo-index-3r0q8ns.svc.aped-4627-b74a.pinecone.io"
)

var (
	Embedder      embeddings.Embedder
	PineconeStore pinecone.Store
)

func initPinecone() {
	initEmbedder()
	initPineconeStore()
	// 上传解析文档
	// 搜索文档
	// 添加到用户提示词中
}

func initEmbedder() {
	embeddingModel, err := openai.New(
		openai.WithBaseURL(EmbeddingBaseUrl),
		openai.WithEmbeddingModel(EmbeddingModel),
		openai.WithToken(AliBaiLianApiKey),
	)
	handleError("初始化embeddingModel失败！", err)

	Embedder, err = embeddings.NewEmbedder(embeddingModel)
	handleError("初始化Embedder失败！", err)

	handleSuccess("初始化Embedder成功！")
}

func initPineconeStore() {
	var err error
	PineconeStore, err = pinecone.New(
		pinecone.WithNameSpace(PineconeNamespace),
		pinecone.WithHost(PineconeHost),
		pinecone.WithAPIKey(PineconeApiKey),
		pinecone.WithEmbedder(Embedder),
	)
	handleError("初始化PineconeStore失败！", err)

	handleSuccess("初始化PineconeStore成功！")
}
