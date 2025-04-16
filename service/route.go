package service

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model/entity"
	"strconv"
)

// CreateRoute 创建线路
func CreateRoute(c *gin.Context) {
	name := c.PostForm("name")
	routeType := c.PostForm("type")
	typeInt, err := strconv.Atoi(routeType)
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	description := c.PostForm("description")
	pointsStr := c.PostForm("points")
	distance := c.PostForm("distance")
	distanceFloat, err := strconv.ParseFloat(distance, 64)
	if name == "" || routeType == "" || status == "" || pointsStr == "" || distance == "" || err != nil {
		fmt.Println(name, routeType, status, description, pointsStr, distance, err)
		common.ErrorResponse(c, common.ParamError)
		return
	}
	var points []entity.GeoPoint
	if err := json.Unmarshal([]byte(pointsStr), &points); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	if len(points) < 2 {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	// 使用雪花算法生成16位的线路ID
	node, err := snowflake.NewNode(1)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	snowflakeID := node.Generate()
	routeId := snowflakeID.String()[:16] // 截取前16位
	route := &entity.Route{
		RouteID:     routeId,
		Name:        name,
		Type:        entity.RouteType(typeInt),
		Status:      entity.RouteStatus(statusInt),
		Description: description,
		Points:      points,
		Distance:    distanceFloat,
	}
	err = entity.InsertRoute(route)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// UpdateRoute 更新线路信息
func UpdateRoute(c *gin.Context) {
	routeId := c.PostForm("routeId")
	name := c.PostForm("name")
	routeType := c.PostForm("type")
	typeInt, err := strconv.Atoi(routeType)
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	description := c.PostForm("description")
	pointsStr := c.PostForm("points")
	distance := c.PostForm("distance")
	distanceFloat, err := strconv.ParseFloat(distance, 64)
	if routeId == "" || name == "" || routeType == "" || typeInt == 0 || statusInt == 0 || description == "" || pointsStr == "" || distance == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	var points []entity.GeoPoint
	if err := json.Unmarshal([]byte(pointsStr), &points); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	if len(points) < 2 {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	route := &entity.Route{
		RouteID:     routeId,
		Name:        name,
		Type:        entity.RouteType(typeInt),
		Status:      entity.RouteStatus(statusInt),
		Description: description,
		Points:      points,
		Distance:    distanceFloat,
	}
	err = entity.UpdateRoute(route)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// GetRouteList 获取线路列表
func GetRouteList(c *gin.Context) {
	var dto entity.FindRouteListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	routes, err := entity.GetRouteList(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, routes)
}

// DeleteRoute 删除线路
func DeleteRoute(c *gin.Context) {
	routeId := c.Query("routeId")
	if routeId == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	err := entity.DeleteRoute(routeId)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// GetRouteTotalCount 获取线路总数
func GetRouteTotalCount(c *gin.Context) {
	totalCount, err := entity.GetRouteTotalCount()
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, totalCount)
}
