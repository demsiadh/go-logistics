package entity

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_logistics/common"
	"go_logistics/config"
	"go_logistics/util"
)

var RouteCollection = config.MongoClient.Database("logistics").Collection("route")

// RouteStatus 线路状态的枚举
type RouteStatus int

const (
	RouteStatusActive   RouteStatus = 1 // 活动中
	RouteStatusInactive RouteStatus = 2 // 不活跃
)

func (s RouteStatus) String() string {
	return [...]string{"启用", "禁用"}[s-1]
}

type RouteType int

const (
	RouteTypeNormal  RouteType = 1 // 普通线路
	RouteTypeQuick   RouteType = 2 // 快速线路
	RouteTypeSpecial RouteType = 3 // 特殊线路
)

func (s RouteType) String() string {
	return [...]string{"常规路线", "快速路线", "特殊路线"}[s-1]
}

// Route 表示一个线路，包含各种属性。
type Route struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RouteID     string             `bson:"routeId" json:"routeId"`
	Name        string             `bson:"name" json:"name"`
	Type        RouteType          `bson:"type" json:"type"`
	Status      RouteStatus        `bson:"status" json:"status"`
	Description string             `bson:"description" json:"description"`
	Points      []common.GeoPoint  `bson:"points" json:"points"`     // 线路点位的坐标，采用 GeoJSON 格式
	Distance    float64            `bson:"distance" json:"distance"` // 线路总里程，单位为公里
	StartOutlet string             `bson:"startOutlet" json:"startOutlet"`
	EndOutlet   string             `bson:"endOutlet" json:"endOutlet"`
	CreateTime  primitive.DateTime `bson:"createTime" json:"-"`
	UpdateTime  primitive.DateTime `bson:"updateTime" json:"-"`
}

// FindRouteListDTO 查询线路列表的参数
type FindRouteListDTO struct {
	RouteID string      `json:"routeId"`
	Name    string      `json:"name"`
	Status  RouteStatus `json:"status"`
	Type    RouteType   `json:"type"`
	Page    common.Page `json:"page"`
}

func (dto *FindRouteListDTO) String() string {
	return fmt.Sprintf("routeId: %s, name: %s, status: %s, type: %s, page: %s", dto.RouteID, dto.Name, dto.Status, dto.Type, dto.Page.String())
}

// InsertRoute 新建线路
func InsertRoute(route *Route) error {
	// 填充时间
	route.CreateTime = util.GetMongoTimeNow()
	route.UpdateTime = util.GetMongoTimeNow()
	_, err := RouteCollection.InsertOne(context.Background(), route)
	return err
}

// GetRouteById 根据线路ID获取线路信息
func GetRouteById(routeId string) (route *Route, err error) {
	filter := bson.M{"routeId": routeId}
	err = RouteCollection.FindOne(context.Background(), filter).Decode(&route)
	return route, err
}

// UpdateRoute 修改线路信息
func UpdateRoute(route *Route) error {
	// 构建过滤条件
	filter := bson.M{"routeId": route.RouteID}
	update := bson.M{
		"$set": bson.M{
			"name":        route.Name,
			"type":        route.Type,
			"description": route.Description,
			"points":      route.Points,
			"distance":    route.Distance,
			"status":      route.Status,
			"startOutlet": route.StartOutlet,
			"endOutlet":   route.EndOutlet,
			"updateTime":  util.GetMongoTimeNow(),
		},
	}
	_, err := RouteCollection.UpdateOne(context.Background(), filter, update)
	return err
}

// DeleteRoute 删除线路
func DeleteRoute(routeId string) error {
	vehicles, err := GetVehicleList(FindVehicleListDTO{
		RouteID: routeId,
	})
	if err != nil {
		return err
	}
	if len(vehicles) > 0 {
		return fmt.Errorf("当前线路存在关联车辆")
	}

	// 构建过滤条件
	filter := bson.M{"routeId": routeId}
	// 执行删除操作
	_, err = RouteCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

// GetRouteList 根据条件查询线路列表
func GetRouteList(dto FindRouteListDTO) (routes []*Route, err error) {
	filter := bson.M{}
	if dto.RouteID != "" {
		filter["routeId"] = bson.M{"$regex": dto.RouteID, "$options": "i"}
	}
	if dto.Name != "" {
		filter["name"] = bson.M{"$regex": dto.Name, "$options": "i"}
	}
	if dto.Type != 0 {
		filter["type"] = dto.Type
	}
	if dto.Status != 0 {
		filter["status"] = dto.Status
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((dto.Page.Skip - 1) * dto.Page.Limit))
	findOptions.SetLimit(int64(dto.Page.Limit))
	findOptions.SetSort(bson.M{"name": 1})

	cursor, err := RouteCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var route Route
		if err := cursor.Decode(&route); err != nil {
			return nil, err
		}
		routes = append(routes, &route)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return routes, nil
}

// GetRouteTotalCount 获取线路总数
func GetRouteTotalCount() (count int64, err error) {
	documents, err := RouteCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return
	}
	return documents, nil
}
