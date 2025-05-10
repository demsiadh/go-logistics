package service

import (
	"fmt"
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
	currentLoad := c.PostForm("currentLoad")
	currentLoadFloat, err := strconv.ParseFloat(currentLoad, 64)
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	routeId := c.PostForm("routeId")
	remarks := c.PostForm("remarks")
	lng := c.PostForm("lng")
	lat := c.PostForm("lat")
	if plateNumber == "" || vType == "" || loadCapacity == "" ||
		currentLoad == "" || status == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}

	vehicle := &entity.Vehicle{
		PlateNumber:  plateNumber,
		Type:         entity.VehicleType(vTypeInt),
		LoadCapacity: loadCapacityFloat,
		CurrentLoad:  currentLoadFloat,
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
	var dto entity.FindVehicleListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	totalCount, err := entity.GetVehicleTotalCount(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, totalCount)
}

// FindMaxRemainingCapacityVehicle 寻找车辆列表中当前剩余载重量最大的车辆
func FindMaxRemainingCapacityVehicle(vehicles []*entity.Vehicle) (*entity.Vehicle, error) {
	if len(vehicles) == 0 {
		return nil, fmt.Errorf("no vehicles available")
	}

	var maxRemaining float64 = -1
	var selectedVehicle *entity.Vehicle

	for _, v := range vehicles {
		remaining := v.LoadCapacity - v.CurrentLoad
		if remaining > maxRemaining {
			maxRemaining = remaining
			selectedVehicle = v
		}
	}

	if selectedVehicle == nil {
		return nil, fmt.Errorf("unable to find vehicle with max remaining capacity")
	}

	return selectedVehicle, nil
}
