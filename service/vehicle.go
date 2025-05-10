package service

import (
	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model/entity"
	"go_logistics/model/vo"
	"strconv"
)

// CreateVehicle 创建车辆
func CreateVehicle(c *gin.Context) {
	plateNumber := c.PostForm("plateNumber")
	vType := c.PostForm("type")
	vTypeInt, err := strconv.Atoi(vType)
	loadCapacity := c.PostForm("loadCapacity")
	loadCapacityFloat, err := strconv.ParseFloat(loadCapacity, 64)
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	routeId := c.PostForm("routeId")
	remarks := c.PostForm("remarks")
	lng := c.PostForm("lng")
	lat := c.PostForm("lat")
	if plateNumber == "" || vType == "" || loadCapacity == "" || status == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}

	vehicle := &entity.Vehicle{
		PlateNumber:  plateNumber,
		Type:         entity.VehicleType(vTypeInt),
		LoadCapacity: loadCapacityFloat,
		Status:       entity.VehicleStatus(statusInt),
		RouteID:      routeId,
		Remarks:      remarks,
		Lng:          lng,
		Lat:          lat,
	}

	err = entity.InsertVehicle(vehicle)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// GetVehicleList 获取车辆列表
func GetVehicleList(c *gin.Context) {
	var dto entity.FindVehicleListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	vehicles, err := entity.GetVehicleList(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	vehicleVOs, err := vo.ToVehicleVOList(vehicles)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, vehicleVOs)
}

// UpdateVehicle 更新车辆信息
func UpdateVehicle(c *gin.Context) {
	plateNumber := c.PostForm("plateNumber")
	vType := c.PostForm("type")
	vTypeInt, err := strconv.Atoi(vType)
	loadCapacity := c.PostForm("loadCapacity")
	loadCapacityFloat, err := strconv.ParseFloat(loadCapacity, 64)
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	routeId := c.PostForm("routeId")
	remarks := c.PostForm("remarks")
	lng := c.PostForm("lng")
	lat := c.PostForm("lat")
	if plateNumber == "" || vType == "" || loadCapacity == "" ||
		status == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}

	vehicle := &entity.Vehicle{
		PlateNumber:  plateNumber,
		Type:         entity.VehicleType(vTypeInt),
		LoadCapacity: loadCapacityFloat,
		Status:       entity.VehicleStatus(statusInt),
		RouteID:      routeId,
		Remarks:      remarks,
		Lng:          lng,
		Lat:          lat,
	}

	err = entity.UpdateVehicle(vehicle)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// DeleteVehicle 删除车辆
func DeleteVehicle(c *gin.Context) {
	id := c.Query("plateNumber")
	if id == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	err := entity.DeleteVehicle(id)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// GetVehicleTotalCount 获取车辆总数
func GetVehicleTotalCount(c *gin.Context) {
	totalCount, err := entity.GetVehicleTotalCount()
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, totalCount)
}
