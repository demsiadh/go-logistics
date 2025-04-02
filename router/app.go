package router

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go_logistics/config"
	"time"
)

func Router() (server *gin.Engine) {
	server = gin.New()
	server.Use(ginzap.Ginzap(config.Log, time.RFC3339, true))
	server.Use(ginzap.RecoveryWithZap(config.Log, true))
	server.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
	return
}
