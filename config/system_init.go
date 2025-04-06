package config

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	MongoClient *mongo.Client
	Log         *zap.Logger
)

func init() {
	initEnvConfig()
	initLog()
	initMongoDB()
}
