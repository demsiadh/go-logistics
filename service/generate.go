package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model/entity"
	"go_logistics/util"
	"math/rand"
)

func GenerateVehicles(c *gin.Context) {
	routes, _ := entity.GetRouteList(entity.FindRouteListDTO{
		Page: common.Page{
			Skip:  1,
			Limit: 1000,
		},
	})
	for i := 0; i < 100; i++ {
		// 随机选择一个线路
		routeIndex := rand.Intn(len(routes))
		selectedRoute := routes[routeIndex]
		vehicle := &entity.Vehicle{
			PlateNumber:  fmt.Sprintf("京A%04d", rand.Intn(10000)),
			Type:         entity.VehicleType(rand.Intn(3) + 1),   // 随机选择车辆类型
			LoadCapacity: float64(rand.Intn(5) + 1),              // 随机选择载重
			Status:       entity.VehicleStatus(rand.Intn(3) + 1), // 随机选择状态
			RouteID:      selectedRoute.RouteID,
			RouteName:    selectedRoute.Name,
			Remarks:      fmt.Sprintf("Test Vehicle %d", i),
			CreateTime:   util.GetMongoTimeNow(),
			UpdateTime:   util.GetMongoTimeNow(),
		}
		_ = entity.InsertVehicle(vehicle)
	}
	common.SuccessResponse(c)
}
