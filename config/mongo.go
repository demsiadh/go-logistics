package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

var MongoClient *mongo.Client

// MongoDBConfig MongoDB配置结构
type MongoDBConfig struct {
	URI        string
	Database   string
	Username   string
	Password   string
	AuthSource string
}

// 初始化MongoDB连接
func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg := getDefaultConfig()
	// 创建客户端选项
	clientOptions := options.Client().ApplyURI(cfg.URI).
		SetAuth(options.Credential{
			Username:   cfg.Username,
			Password:   cfg.Password,
			AuthSource: cfg.AuthSource,
		}).
		SetMaxPoolSize(100).                 // 连接池大小
		SetMinPoolSize(10).                  // 最小连接数
		SetMaxConnIdleTime(30 * time.Minute) // 最大空闲时间

	// 建立连接
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		Log.Error("MongoDB连接失败: ", zap.Any(ERROR, err))
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		Log.Error("MongoDB心跳检测失败: ", zap.Any(ERROR, err))
	}

	MongoClient = client
	Log.Info("成功连接到MongoDB!")
}

func getDefaultConfig() (cfg *MongoDBConfig) {
	cfg = &MongoDBConfig{
		URI:        "mongodb://logistics:logistics@tccloud:27017",
		Database:   "logistics",
		Username:   "logistics",
		Password:   "logistics",
		AuthSource: "logistics",
	}
	return
}
