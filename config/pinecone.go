package config

import (
	pineconeClient "github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores/pinecone"
)

const (
	EmbeddingModel    = "text-embedding-v3"
	EmbeddingBaseUrl  = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	PineconeNamespace = "logistics"
	PineconeHost      = "logistics-3r0q8ns.svc.aped-4627-b74a.pinecone.io"
)

var (
	Embedder            embeddings.Embedder
	PineconeStore       pinecone.Store
	PineconeIndexClient *pineconeClient.IndexConnection
)

func initPinecone() {
	initEmbedder()
	initPineconeStore()
	initPineconeIndexClient()
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

func initPineconeIndexClient() {
	pcClient, err := pineconeClient.NewClient(pineconeClient.NewClientParams{
		ApiKey: PineconeApiKey,
	})
	handleError("初始化PineconeClient失败！", err)
	PineconeIndexClient, err = pcClient.IndexWithNamespace(PineconeHost, PineconeNamespace)
	handleError("初始化PineconeIndexClient失败！", err)
	handleSuccess("初始化PineconeIndexClient成功！")
}
