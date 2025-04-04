package router

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go_logistics/config"
	"go_logistics/service"
	"time"
)

func Router() (server *gin.Engine) {
	server = gin.New()
	server.Use(ginzap.Ginzap(config.Log, time.RFC3339, true))
	server.Use(ginzap.RecoveryWithZap(config.Log, true))
	server.GET("/user/getByName", service.GetUserByName)
	server.POST("/user/create", service.CreateUser)
	server.GET("/user/list", service.GetUserList)
	server.PUT("/user/update", service.UpdateUser)
	server.DELETE("/user/delete", service.DeleteUser)
	return
}
