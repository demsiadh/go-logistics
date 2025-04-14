package service

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go_logistics/common"
	"go_logistics/model/entity"
	"strconv"
)

const baiduMapAK = "SBG5d7X8VNM0nlfBfIjRXWvJ01LAR3Bg" // 替换为你的百度地图API Key

// CreateRoute 创建线路
func CreateRoute(c *gin.Context) {
	name := c.PostForm("name")
	routeType := c.PostForm("type")
	typeInt, err := strconv.Atoi(routeType)
	status := c.PostForm("status")
	statusInt, err := strconv.Atoi(status)
	description := c.PostForm("description")
	pointsStr := c.PostForm("points")
	if name == "" || routeType == "" || typeInt == 0 || statusInt == 0 || description == "" || pointsStr == "" || err != nil {
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
	distance, err := calculateDistance(points)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
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
		Distance:    distance,
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
	if routeId == "" || name == "" || routeType == "" || typeInt == 0 || statusInt == 0 || description == "" || pointsStr == "" || err != nil {
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
	distance, err := calculateDistance(points)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	route := &entity.Route{
		RouteID:     routeId,
		Name:        name,
		Type:        entity.RouteType(typeInt),
		Status:      entity.RouteStatus(statusInt),
		Description: description,
		Points:      points,
		Distance:    distance,
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

// calculateDistance 计算点位之间的距离
func calculateDistance(points []entity.GeoPoint) (float64, error) {
	client := resty.New()
	var origins, destinations string
	for i, point := range points {
		if i > 0 {
			origins += "|"
			destinations += "|"
		}
		origins += fmt.Sprintf("%f,%f", point.Coordinates[1], point.Coordinates[0])
		destinations += fmt.Sprintf("%f,%f", point.Coordinates[1], point.Coordinates[0])
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"origins":      origins,
			"destinations": destinations,
			"ak":           baiduMapAK,
		}).
		Get("http://api.map.baidu.com/routematrix/v2/driving")
	if err != nil {
		return 0, err
	}

	var result struct {
		Result []struct {
			Distance struct {
				Value int `json:"value"`
			} `json:"distance"`
		} `json:"result"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return 0, err
	}

	var totalDistance float64
	for _, r := range result.Result {
		totalDistance += float64(r.Distance.Value) / 1000 // 转换为公里
	}

	return totalDistance, nil
}
