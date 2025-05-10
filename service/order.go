package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go_logistics/config"
	"go_logistics/model/vo"
	"go_logistics/util"
	"math"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model/entity"
	"strconv"
)

// generateOrderID 生成订单ID
func generateOrderID() (string, error) {
	currentDate := time.Now().Format("20060102")
	count, err := entity.GetOrderCountByDate(currentDate)
	if err != nil {
		return "", err
	}
	count++
	orderID := fmt.Sprintf("LY%s%04d", currentDate, count)
	return orderID, nil
}

// CreateOrder 创建订单
func CreateOrder(c *gin.Context) {
	customerName := c.PostForm("customerName")
	phone := c.PostForm("phone")
	startAddress := c.PostForm("startAddress")
	startLng := c.PostForm("startLng")
	startLat := c.PostForm("startLat")
	endAddress := c.PostForm("endAddress")
	endLng := c.PostForm("endLng")
	endLat := c.PostForm("endLat")
	remark := c.PostForm("remark")
	weightStr := c.PostForm("weight")
	weight, err := strconv.ParseFloat(weightStr, 64)
	if customerName == "" || phone == "" || startAddress == "" || endAddress == "" ||
		startLng == "" || startLat == "" || endLng == "" || endLat == "" || weightStr == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	orderID, err := generateOrderID()
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	order := &entity.Order{
		OrderID:      orderID,
		CustomerName: customerName,
		Phone:        phone,
		Status:       entity.Pending,
		StartAddress: startAddress,
		StartLng:     startLng,
		StartLat:     startLat,
		EndAddress:   endAddress,
		EndLng:       endLng,
		EndLat:       endLat,
		Weight:       weight,
		Remark:       remark,
	}
	err = entity.InsertOrder(order)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			config.Log.Warn("协程超时或被取消", zap.Error(ctx.Err()))
			return
		default:
			completeOrder(orderID)
		}
	}(ctx)

	common.SuccessResponse(c)
}

// GetOrderList 获取订单列表
func GetOrderList(c *gin.Context) {
	var dto entity.FindOrderListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	orders, err := entity.GetOrderList(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, orders)
}

// UpdateOrder 更新订单信息
func UpdateOrder(c *gin.Context) {
	orderId := c.PostForm("orderId")
	customerName := c.PostForm("customerName")
	phone := c.PostForm("phone")
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	remark := c.PostForm("remark")
	if orderId == "" || customerName == "" || phone == "" || status == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	order := &entity.Order{
		OrderID:      orderId,
		CustomerName: customerName,
		Phone:        phone,
		Status:       entity.OrderStatus(statusInt),
		Remark:       remark,
	}
	err = entity.UpdateOrder(order)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// DeleteOrder 删除订单
func DeleteOrder(c *gin.Context) {
	orderId := c.Query("orderId")
	if orderId == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	err := entity.DeleteOrder(orderId)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// GetOrderTotalCount 获取订单总数
func GetOrderTotalCount(c *gin.Context) {
	var dto entity.FindOrderListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	totalCount, err := entity.GetOrderTotalCount(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, totalCount)
}

func GetOrderVO(c *gin.Context) {
	orderId := c.Query("orderId")
	if orderId == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	order, err := entity.GetOrderById(orderId)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	orderVO, err := vo.ToOrderVO(order)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, orderVO)
}

// LockPool Key-based Lock Pool
var LockPool = struct {
	mu sync.Mutex
	m  map[string]*sync.Mutex
}{m: make(map[string]*sync.Mutex)}

func GetOrderLock(orderID string) *sync.Mutex {
	LockPool.mu.Lock()
	defer LockPool.mu.Unlock()

	if _, exists := LockPool.m[orderID]; !exists {
		LockPool.m[orderID] = &sync.Mutex{}
	}
	return LockPool.m[orderID]
}

func GetVehicleLock(plateNumber string) *sync.Mutex {
	LockPool.mu.Lock()
	defer LockPool.mu.Unlock()

	if _, exists := LockPool.m[plateNumber]; !exists {
		LockPool.m[plateNumber] = &sync.Mutex{}
	}
	return LockPool.m[plateNumber]
}

func completeOrder(orderId string) {
	orderMu := GetOrderLock(orderId)
	orderMu.Lock()
	defer orderMu.Unlock()

	// 获取订单信息
	order, err := entity.GetOrderById(orderId)
	if err != nil {
		config.Log.Warn("填充订单信息失败！", zap.String("orderId", orderId))
		return
	}

	// 查找起点网点
	startOutlet, err := findNearOutlet(order.StartLng, order.StartLat)
	if err != nil {
		config.Log.Warn("查询起点网点失败！", zap.String("orderId", orderId))
		return
	}

	// 判断是否在范围内...
	isInScope, err := util.IsPointInGeoPointSlice(order.StartLng, order.StartLat, startOutlet.Scope)
	if err != nil || !isInScope {
		config.Log.Warn("起点不在网点营业范围内！", zap.String("orderId", orderId))
		return
	}

	// 查找终点网点...
	endOutlet, err := findNearOutlet(order.EndLng, order.EndLat)
	if err != nil {
		config.Log.Warn("查询终点网点失败！", zap.String("orderId", orderId))
		return
	}
	// 判断是否在范围内...
	isInScope, err = util.IsPointInGeoPointSlice(order.StartLng, order.StartLat, startOutlet.Scope)
	if err != nil || !isInScope {
		config.Log.Warn("终点不在网点营业范围内！", zap.String("orderId", orderId))
		return
	}

	// 查询线路与车辆...
	routes, err := entity.GetRouteByOutlets(startOutlet.ID.Hex(), endOutlet.ID.Hex())
	if err != nil {
		config.Log.Warn("获取线路失败！", zap.String("orderId", orderId))
		return
	}

	var vehicles []*entity.Vehicle
	for _, route := range routes {
		tempVehicles, err := entity.GetVehicleByRouteId(route.RouteID)
		if err != nil {
			config.Log.Warn("获取车辆失败！", zap.String("orderId", orderId))
			return
		}
		if len(tempVehicles) > 0 {
			vehicles = append(vehicles, tempVehicles...)
		}
	}

	vehicle, err := FindMaxRemainingCapacityVehicle(vehicles)
	if err != nil {
		config.Log.Warn("获取最优车辆失败！", zap.String("orderId", orderId))
		return
	}

	// 获取车辆锁
	vehicleMu := GetVehicleLock(vehicle.PlateNumber)
	vehicleMu.Lock()
	defer vehicleMu.Unlock()

	// 检查载重
	if vehicle.CurrentLoad+order.Weight > vehicle.LoadCapacity {
		config.Log.Warn("当前车辆超载！", zap.String("orderId", orderId))
		vehicle.Status = entity.InTransit
		err = entity.UpdateVehicle(vehicle)
		if err != nil {
			config.Log.Warn("更新车辆状态失败！", zap.String("orderId", orderId))
		}
		return
	}

	// 更新车辆负载
	vehicle.CurrentLoad += order.Weight
	err = entity.UpdateVehicle(vehicle)
	if err != nil {
		config.Log.Warn("更新车辆状态失败！", zap.String("orderId", orderId))
		return
	}

	// 更新订单状态
	order.StartOutletId = startOutlet.ID.Hex()
	order.EndOutletId = endOutlet.ID.Hex()
	order.TransPortVehicle = vehicle.PlateNumber
	err = entity.CompleteOrder(order)
	if err != nil {
		config.Log.Warn("更新订单状态失败！", zap.String("orderId", orderId))
		return
	}
}

// 查找最近的网点
func findNearOutlet(lng string, lat string) (*entity.Outlet, error) {
	outlets, err := entity.GetOutletList(entity.FindOutletListDTO{
		Page: common.Page{
			Skip:  1,
			Limit: 1000,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(outlets) == 0 {
		return nil, nil
	}

	var nearest *entity.Outlet
	minDist := math.MaxFloat64

	for i := range outlets {
		outlet := outlets[i]
		distanceMeters, err := util.GetDistanceFromString(lng, lat, outlet.Lng, outlet.Lat)
		if err != nil {
			return nil, err
		}
		if distanceMeters < minDist {
			minDist = distanceMeters
			nearest = outlet
		}
	}

	if nearest == nil {
		return nil, fmt.Errorf("no valid outlet found")
	}

	return nearest, nil
}
