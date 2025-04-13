package service

import (
	"fmt"
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
	address := c.PostForm("address")
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	remark := c.PostForm("remark")
	if customerName == "" || phone == "" || address == "" || status == "" || err != nil {
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
		Address:      address,
		Status:       entity.OrderStatus(statusInt),
		Remark:       remark,
	}
	err = entity.InsertOrder(order)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
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
	address := c.PostForm("address")
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	remark := c.PostForm("remark")
	if orderId == "" || customerName == "" || phone == "" || address == "" || status == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	order := &entity.Order{
		OrderID:      orderId,
		CustomerName: customerName,
		Phone:        phone,
		Address:      address,
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
	totalCount, err := entity.GetOrderTotalCount()
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, totalCount)
}
