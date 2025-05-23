package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model/entity"
	"strconv"
)

// CreateOutlet 创建网点
func CreateOutlet(c *gin.Context) {
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	province := c.PostForm("province")
	city := c.PostForm("city")
	detailAddress := c.PostForm("detailAddress")
	businessHours := c.PostForm("businessHours")
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	remark := c.PostForm("remark")
	lng := c.PostForm("lng")
	lat := c.PostForm("lat")
	scopeStr := c.PostForm("scope")
	if name == "" || phone == "" || province == "" || city == "" || scopeStr == "" || detailAddress == "" ||
		businessHours == "" || status == "" || lng == "" || lat == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}

	var scope []common.GeoPoint
	if err := json.Unmarshal([]byte(scopeStr), &scope); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	if len(scope) <= 2 {
		common.ErrorResponse(c, common.ParamError)
		return
	}

	outlet := &entity.Outlet{
		Name:          name,
		Phone:         phone,
		Province:      province,
		City:          city,
		DetailAddress: detailAddress,
		BusinessHours: businessHours,
		Status:        entity.OutletStatus(statusInt),
		Remark:        remark,
		Lng:           lng,
		Lat:           lat,
		Scope:         scope,
	}
	err = entity.InsertOutlet(outlet)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponse(c)
}

// GetOutletList 获取网点列表
func GetOutletList(c *gin.Context) {
	var dto entity.FindOutletListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	outlets, err := entity.GetOutletList(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponseWithData(c, outlets)
}

// UpdateOutlet 更新网点信息
func UpdateOutlet(c *gin.Context) {
	outletId := c.PostForm("id")
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	detailAddress := c.PostForm("detailAddress")
	businessHours := c.PostForm("businessHours")
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	remark := c.PostForm("remark")
	lng := c.PostForm("lng")
	lat := c.PostForm("lat")
	scopeStr := c.PostForm("scope")
	if name == "" || phone == "" || detailAddress == "" ||
		businessHours == "" || status == "" || lng == "" ||
		lat == "" || scopeStr == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}

	var scope []common.GeoPoint
	if err := json.Unmarshal([]byte(scopeStr), &scope); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	if len(scope) <= 2 {
		common.ErrorResponse(c, common.ParamError)
		return
	}

	// 构建过滤条件
	outlet := &entity.Outlet{
		Name:          name,
		Phone:         phone,
		DetailAddress: detailAddress,
		BusinessHours: businessHours,
		Status:        entity.OutletStatus(statusInt),
		Remark:        remark,
		Lng:           lng,
		Lat:           lat,
		Scope:         scope,
	}
	err = entity.UpdateOutlet(outletId, outlet)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponse(c)
}

// DeleteOutlet 删除网点
func DeleteOutlet(c *gin.Context) {
	outletId := c.Query("outletId")
	if outletId == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	err := entity.DeleteOutlet(outletId)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponse(c)
}

// GetOutletTotalCount 获取网点总数
func GetOutletTotalCount(c *gin.Context) {
	var dto entity.FindOutletListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	totalCount, err := entity.GetOutletTotalCount(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponseWithData(c, totalCount)
}

// GetAllProvincesAndCities 获取所有省份和城市
func GetAllProvincesAndCities(c *gin.Context) {
	result, err := entity.GetAllProvincesAndCities()
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponseWithData(c, result)
}

func GetOutletById(c *gin.Context) {
	outletId := c.Query("outletId")
	if outletId == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	outlet, err := entity.GetOutletById(outletId)
	if err != nil {
		common.ErrorResponse(c, common.ServerError(err.Error()))
		return
	}
	common.SuccessResponseWithData(c, outlet)
}
