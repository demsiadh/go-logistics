package router

import (
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/config"
	"go_logistics/model/entity"
	"go_logistics/service"
	"go_logistics/util"
	"time"
)

func Router() (server *gin.Engine) {
	server = gin.New()
	// 日志中间件
	server.Use(ginzap.Ginzap(config.Log, "2006-01-02 15:04:05.000 CST", false))
	server.Use(ginzap.RecoveryWithZap(config.Log, true))

	// CORS 配置（严格模式）
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 明确指定前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Logistics-Custom-Header", "logistics_token"},
		ExposeHeaders:    []string{"Content-Length", "logistics_token"},
		AllowCredentials: true, // 关键配置
		MaxAge:           12 * time.Hour,
	}))

	// 全局中间件（所有路由都会经过）
	server.Use(TokenAuthMiddleware())

	apiGroup := server.Group("/api")
	userGroup := apiGroup.Group("/user")
	{
		userGroup.GET("/getByName", service.GetUserByName)
		userGroup.POST("/create", service.CreateUser)
		userGroup.POST("/list", service.GetUserList)
		userGroup.PUT("/update", service.UpdateUser)
		userGroup.DELETE("/delete", service.DeleteUser)
		userGroup.POST("/login", service.LoginUser)
		userGroup.GET("/loginStatus", service.GetUserLoginStatus)
		userGroup.POST("/total", service.GetTotalCount)
	}
	orderGroup := apiGroup.Group("/order")
	{
		orderGroup.POST("/create", service.CreateOrder)
		orderGroup.POST("/list", service.GetOrderList)
		orderGroup.POST("/total", service.GetOrderTotalCount)
		orderGroup.PUT("/update", service.UpdateOrder)
		orderGroup.DELETE("/delete", service.DeleteOrder)
	}
	outletGroup := apiGroup.Group("/outlet")
	{
		outletGroup.POST("/create", service.CreateOutlet)
		outletGroup.POST("/list", service.GetOutletList)
		outletGroup.POST("/total", service.GetOutletTotalCount)
		outletGroup.PUT("/update", service.UpdateOutlet)
		outletGroup.DELETE("/delete", service.DeleteOutlet)
		outletGroup.GET("/allProvincesAndCities", service.GetAllProvincesAndCities)
		outletGroup.GET("/id", service.GetOutletById)
	}
	routeGroup := apiGroup.Group("/route")
	{
		routeGroup.POST("/create", service.CreateRoute)
		routeGroup.POST("/list", service.GetRouteList)
		routeGroup.POST("/total", service.GetRouteTotalCount)
		routeGroup.PUT("/update", service.UpdateRoute)
		routeGroup.DELETE("/delete", service.DeleteRoute)
	}

	vehicleGroup := apiGroup.Group("/vehicle")
	{
		vehicleGroup.POST("/create", service.CreateVehicle)
		vehicleGroup.POST("/list", service.GetVehicleList)
		vehicleGroup.POST("/total", service.GetVehicleTotalCount)
		vehicleGroup.PUT("/update", service.UpdateVehicle)
		vehicleGroup.DELETE("/delete", service.DeleteVehicle)
	}
	homeGroup := apiGroup.Group("/home")
	{
		homeGroup.GET("/outlet", service.GetOutletView)
		homeGroup.GET("/order", service.GetOrderView)
		homeGroup.GET("/vehicle", service.GetVehicleView)
		homeGroup.GET("/route", service.GetRouteView)
	}
	generateGroup := apiGroup.Group("/generate")
	{
		generateGroup.GET("/vehicle", service.GenerateVehicles)
	}
	llmGroup := apiGroup.Group("/llm")
	{
		llmGroup.GET("chat", service.ChatLLM)
	}

	return
}

// TokenAuthMiddleware Token 校验中间件
func TokenAuthMiddleware() gin.HandlerFunc {
	// 定义白名单路径（支持精确匹配）
	whitelist := map[string]bool{
		"/api/user/login": true, // 登录接口
		"/api/llm/chat":   true,
	}

	return func(c *gin.Context) {
		currentPath := c.FullPath() // 获取注册的路由路径（非请求URI）
		// 检查白名单
		if _, ok := whitelist[currentPath]; ok {
			c.Next()
			return
		}

		var token string
		token = c.GetHeader("logistics_token")
		if token == "" {
			common.AbortResponse(c, common.NotLogin)
			return
		}

		claims, err := util.CheckToken(token)
		if err != nil {
			common.AbortResponse(c, common.NotLogin)
			return
		}
		if claims == nil {
			common.AbortResponse(c, common.NotLogin)
			return
		}
		name := claims.Name
		user, err := entity.GetUserByName(name)
		if err != nil {
			common.AbortResponse(c, common.RecordNotFound)
			return
		}
		if user.Status == entity.Banned {
			common.AbortResponse(c, common.UserBanned)
			return
		}
		if user.Status == entity.Deleted {
			common.AbortResponse(c, common.UserDeleted)
			return
		}
		c.Set("name", claims.Name)
		c.Next()
		return
	}
}
