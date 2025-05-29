package config

import (
	"fmt"
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

func initPinecone() {
	embedding, err := openai.New(
		openai.WithBaseURL(EmbeddingBaseUrl),
		openai.WithEmbeddingModel(EmbeddingModel),
		openai.WithToken(AliBaiLianApiKey),
	)
	embedder, err := embeddings.NewEmbedder(embedding)
	store, err := pinecone.New(
		pinecone.WithNameSpace(PineconeNamespace),
		pinecone.WithHost(PineconeHost),
		pinecone.WithAPIKey(PineconeApiKey),
		pinecone.WithEmbedder(embedder),
	)
	Log.Info("初始化pinecone成功！")
	fmt.Println(store)

	if err != nil {
		fmt.Println(err)
	}
	// 接入向量模型
	// 上传解析文档
	// 搜索文档
	// 添加到用户提示词中
}
