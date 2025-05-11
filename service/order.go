package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go_logistics/config"
	"go_logistics/model/vo"
	"go_logistics/util"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model/entity"
	"strconv"
)

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
	orderID, err := util.GenerateOrderID()
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
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
		common.ErrorResponse(c, common.ServerError(err.Error()))
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
			completeDataOrder(orderID)
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
		common.ErrorResponse(c, common.ServerError(err.Error()))
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
		common.ErrorResponse(c, common.ServerError(err.Error()))
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
		common.ErrorResponse(c, common.ServerError(err.Error()))
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
		common.ErrorResponse(c, common.ServerError(err.Error()))
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
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	orderVO, err := vo.ToOrderVO(order)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponseWithData(c, orderVO)
}

func completeDataOrder(orderId string) {
	orderMu := util.GetOrderLock(orderId)
	orderMu.Lock()
	defer orderMu.Unlock()

	// 获取订单信息
	order, err := entity.GetOrderById(orderId)
	if err != nil {
		msg := "填充订单信息失败！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
		return
	}

	// 查找起点网点
	startOutlet, err := findNearOutlet(order.StartLng, order.StartLat)
	if err != nil {
		msg := "查询起点网点失败！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
		return
	}

	// 判断是否在范围内...
	isInScope, err := util.IsPointInGeoPointSlice(order.StartLng, order.StartLat, startOutlet.Scope)
	if err != nil || !isInScope {
		msg := "起点不在网点营业范围内！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
		return
	}

	// 查找终点网点...
	endOutlet, err := findNearOutlet(order.EndLng, order.EndLat)
	if err != nil {
		msg := "查询终点网点失败！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
		return
	}

	// 判断是否在范围内...
	isInScope, err = util.IsPointInGeoPointSlice(order.EndLng, order.EndLat, endOutlet.Scope)
	if err != nil || !isInScope {
		msg := "终点不在网点营业范围内！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
		return
	}
	if startOutlet.ID.Hex() == endOutlet.ID.Hex() {
		msg := "起点与终点在同一个网点！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		// 更新订单状态
		order.StartOutletId = startOutlet.ID.Hex()
		order.EndOutletId = endOutlet.ID.Hex()
		order.Remark = ""
		err = entity.CompleteDataOrder(order)
		if err != nil {
			msg := "更新订单状态失败！"
			config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
			_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
			return
		}
		return
	}

	// 查询线路与车辆...
	routes, err := entity.GetRouteByOutlets(startOutlet.ID.Hex(), endOutlet.ID.Hex())
	if err != nil {
		msg := "获取线路失败！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
		return
	}

	var vehicles []*entity.Vehicle
	for _, route := range routes {
		tempVehicles, err := entity.GetVehicleByRouteId(route.RouteID)
		if err != nil {
			msg := "获取车辆失败！"
			config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
			_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
			return
		}
		if len(tempVehicles) > 0 {
			vehicles = append(vehicles, tempVehicles...)
		}
	}

	vehicle, err := FindMaxRemainingCapacityVehicle(vehicles)
	if err != nil {
		msg := "获取最优车辆失败！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
		return
	}

	// 获取车辆锁
	vehicleMu := util.GetVehicleLock(vehicle.PlateNumber)
	vehicleMu.Lock()
	defer vehicleMu.Unlock()

	// 检查载重
	if vehicle.CurrentLoad+order.Weight > vehicle.LoadCapacity {
		msg := "当前车辆超载！"
		config.Log.Warn(msg, zap.String("orderId", orderId))
		vehicle.Status = entity.InTransit
		err = entity.UpdateVehicle(vehicle)
		if err != nil {
			_ = updateOrderRemark(order, msg+" 更新车辆状态失败: "+err.Error())
		} else {
			_ = updateOrderRemark(order, msg)
		}
		return
	}

	// 定义精度常量（保留4位小数，即精确到0.1公斤）
	const precisionFactor = 10000 // 10^4 = 10000

	// 更新车辆负载
	vehicle.CurrentLoad = roundToPrecision(vehicle.CurrentLoad+order.Weight, precisionFactor)
	err = entity.UpdateVehicle(vehicle)
	if err != nil {
		msg := "更新车辆状态失败！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
		return
	}

	// 更新订单状态
	order.StartOutletId = startOutlet.ID.Hex()
	order.EndOutletId = endOutlet.ID.Hex()
	order.Remark = ""
	order.TransPortVehicle = vehicle.PlateNumber
	err = entity.CompleteDataOrder(order)
	if err != nil {
		msg := "更新订单状态失败！"
		config.Log.Warn(msg, zap.String("orderId", orderId), zap.Error(err))
		_ = updateOrderRemark(order, msg+" 错误原因: "+err.Error())
		return
	}
}

// 四舍五入到指定精度
func roundToPrecision(value float64, factor int) float64 {
	return float64(int64(value*float64(factor)+0.5)) / float64(factor)
}

// updateOrderRemark 更新订单的备注字段
func updateOrderRemark(order *entity.Order, remark string) error {
	order.Remark = remark
	return entity.UpdateOrder(order)
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

// DispatchOrder 手动调度订单
func DispatchOrder(c *gin.Context) {
	orderId := c.Query("orderId")
	if orderId == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	completeDataOrder(orderId)
	common.SuccessResponse(c)
}
