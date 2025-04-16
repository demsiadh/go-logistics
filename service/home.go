package service

import (
	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model/entity"
	"time"
)

func GetOutletView(c *gin.Context) {
	outlets, err := entity.GetOutletList(entity.FindOutletListDTO{
		Page: common.Page{
			Skip:  1,
			Limit: 10000,
		},
	})
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
	}
	common.SuccessResponseWithData(c, outlets)
}

func GetOrderView(c *gin.Context) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)       // 昨天
	startTime := yesterday.AddDate(0, 0, -6) // 从昨天往前推七天

	dto := entity.FindOrderListDTO{
		StartTime: startTime,
		EndTime:   yesterday,
		Page: common.Page{
			Skip:  1,
			Limit: 10000,
		},
	}

	orders, err := entity.GetOrderList(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
	}
	result := GroupOrdersByDateAndStatus(orders)
	common.SuccessResponseWithData(c, result)
}

// GroupOrdersByDateAndStatus 将订单列表按日期和状态分组
func GroupOrdersByDateAndStatus(orders []*entity.Order) map[string]any {
	result := make(map[string]any)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)       // 昨天
	startTime := yesterday.AddDate(0, 0, -6) // 从昨天往前推七天
	endTime := yesterday                     // 修复：定义 endTime 为 yesterday

	// 初始化所有日期的条目
	currentDate := startTime
	for currentDate.Before(endTime.AddDate(0, 0, 1)) { // 使用 endTime
		dateStr := currentDate.Format("2006-01-02")
		result[dateStr] = map[string]int{
			"total":      0,
			"processing": 0,
			"completed":  0,
			"pending":    0,
			"canceled":   0,
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	for _, order := range orders {
		// 使用 ToTime() 方法将 primitive.DateTime 转换为 time.Time
		createTime := order.CreateTime.Time()
		date := createTime.Format("2006-01-02") // 格式化日期

		// 更新统计
		stats := result[date].(map[string]int)
		stats["total"]++
		switch order.Status {
		case entity.OrderStatus(1):
			stats["pending"]++
		case entity.OrderStatus(2):
			stats["processing"]++
		case entity.OrderStatus(3):
			stats["completed"]++
		case entity.OrderStatus(4):
			stats["canceled"]++
		}
	}

	return result
}

func GetVehicleView(c *gin.Context) {
	dto := entity.FindVehicleListDTO{
		Page: common.Page{
			Skip:  1,
			Limit: 10000,
		},
	}

	vehicles, err := entity.GetVehicleList(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
	}
	result := GroupVehiclesByStatus(vehicles)
	common.SuccessResponseWithData(c, result)
}

// GroupVehiclesByStatus 将车辆列表按状态分组
func GroupVehiclesByStatus(vehicles []*entity.Vehicle) map[string]int {
	result := map[string]int{
		"InTransit":   0,
		"Maintenance": 0,
		"Free":        0,
	}

	for _, vehicle := range vehicles {
		switch vehicle.Status {
		case entity.VehicleStatus(1):
			result["InTransit"]++
		case entity.VehicleStatus(2):
			result["Maintenance"]++
		case entity.VehicleStatus(3):
			result["Free"]++
		default:
			// 处理未知状态
			result["unknown"]++
		}
	}

	return result
}

func GetRouteView(c *gin.Context) {
	routes, err := entity.GetRouteList(entity.FindRouteListDTO{
		Page: common.Page{
			Skip:  1,
			Limit: 10000,
		},
	})
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
	}
	result := GroupRoutesByType(routes)
	common.SuccessResponseWithData(c, result)
}

// GroupRoutesByType 将线路列表按类型分组并统计里程
func GroupRoutesByType(routes []*entity.Route) map[string]float64 {
	result := map[string]float64{
		"Normal":  0,
		"Quick":   0,
		"special": 0,
		"unknown": 0,
	}

	for _, route := range routes {
		switch route.Type {
		case entity.RouteType(1):
			result["Normal"] += route.Distance
		case entity.RouteType(2):
			result["Quick"] += route.Distance
		case entity.RouteType(3):
			result["special"] += route.Distance
		default:
			result["unknown"] += route.Distance
		}
	}

	return result
}
