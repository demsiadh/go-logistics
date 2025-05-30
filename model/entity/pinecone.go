package entity

import (
	"context"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	"go.uber.org/zap"
	"go_logistics/config"
	"io"
	"mime/multipart"
)

const (
	ChunkSize    = 256
	ChunkOverlap = 50
)

func InsertVector(ctx context.Context, file multipart.File) (ids []string, err error) {
	_, err = file.Seek(io.SeekStart, io.SeekStart)
	if err != nil {
		config.Log.Error("文件 seek 失败！", zap.Error(err))
		return
	}
	loader := documentloaders.NewText(file)
	docs, err := loader.LoadAndSplit(ctx, textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(ChunkSize),
		textsplitter.WithChunkOverlap(ChunkOverlap),
	))
	if err != nil {
		return
	}

	ids, err = config.PineconeStore.AddDocuments(ctx, docs)
	if err != nil {
		return
	}
	return
}

func SearchAndExtractContents(ctx context.Context, query string, topK int) (textSegments []string, err error) {
	// 1. 向量数据库相似性搜索
	docs, err := config.PineconeStore.SimilaritySearch(ctx, query, topK)
	if err != nil {
		return
	}

	// 2. 检查是否找到结果
	if len(docs) == 0 {
		return
	}

	// 3. 提取每个文档的 Content 字段
	for _, doc := range docs {
		textSegments = append(textSegments, doc.PageContent)
	}

	return
}

func DeleteVector(ctx context.Context, ids []string) (err error) {
	err = config.PineconeIndexClient.DeleteVectorsById(&ctx, ids)
	return
}
