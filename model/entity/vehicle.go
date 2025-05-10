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

var VehicleCollection = config.MongoClient.Database("logistics").Collection("vehicle")

// VehicleStatus 车辆状态
type VehicleStatus int

const (
	InTransit   VehicleStatus = 1
	Maintenance VehicleStatus = 2
	Free        VehicleStatus = 3
)

func (s VehicleStatus) String() string {
	textMap := map[VehicleStatus]string{
		InTransit:   "运行中",
		Maintenance: "维修中",
		Free:        "空闲",
	}
	return textMap[s]
}

type VehicleType int

const (
	Truck   VehicleType = 1
	Minibus VehicleType = 2
	Pickup  VehicleType = 3
)

func (s VehicleType) String() string {
	textMap := map[VehicleType]string{
		Truck:   "货车",
		Minibus: "面包车",
		Pickup:  "皮卡",
	}
	return textMap[s]
}

// Vehicle 车辆结构
type Vehicle struct {
	ID           string             `bson:"_id,omitempty" json:"id"`
	PlateNumber  string             `bson:"plateNumber" json:"plateNumber"`
	Type         VehicleType        `bson:"type" json:"type"`
	LoadCapacity float64            `bson:"loadCapacity" json:"loadCapacity"`
	CurrentLoad  float64            `bson:"currentLoad" json:"currentLoad"`
	Status       VehicleStatus      `bson:"status" json:"status"`
	RouteID      string             `bson:"routeId" json:"routeId"`
	RouteName    string             `bson:"routeName" json:"routeName"`
	Remarks      string             `bson:"remarks" json:"remarks"`
	Lng          string             `bson:"lng" json:"lng"`
	Lat          string             `bson:"lat" json:"lat"`
	CreateTime   primitive.DateTime `bson:"createTime" json:"-"`
	UpdateTime   primitive.DateTime `bson:"updateTime" json:"-"`
}

// FindVehicleListDTO 查询车辆列表的参数
type FindVehicleListDTO struct {
	PlateNumber string        `json:"plateNumber"`
	Type        VehicleType   `json:"type"`
	Status      VehicleStatus `json:"status"`
	RouteID     string        `json:"routeId"`
	RouteName   string        `json:"routeName"`
	Page        common.Page   `json:"page"`
}

func (dto *FindVehicleListDTO) String() string {
	return fmt.Sprintf("plateNumber: %s, type: %d, status: %d, routeId: %s, routeName: %s, page: %s",
		dto.PlateNumber, dto.Type, dto.Status, dto.RouteID, dto.RouteName, dto.Page.String())
}

// InsertVehicle 新建车辆
func InsertVehicle(vehicle *Vehicle) error {
	// 检查车牌号是否已存在
	var existingVehicle Vehicle
	err := VehicleCollection.FindOne(context.Background(), bson.M{"plateNumber": vehicle.PlateNumber}).Decode(&existingVehicle)
	if err == nil {
		return fmt.Errorf("车牌号 %s 已存在", vehicle.PlateNumber)
	}

	// 检查线路是否存在
	var route *Route
	if vehicle.RouteID != "" {
		route, err = GetRouteById(vehicle.RouteID)
		if err != nil {
			return err
		}
		if route == nil {
			return fmt.Errorf("线路不存在")
		}
	} else {
		route = &Route{
			Name: "",
		}
	}
	// 填充时间
	vehicle.CreateTime = util.GetMongoTimeNow()
	vehicle.UpdateTime = util.GetMongoTimeNow()
	// 填充线路名称
	vehicle.RouteName = route.Name

	_, err = VehicleCollection.InsertOne(context.Background(), vehicle)
	return err
}

// UpdateVehicle 修改车辆信息
func UpdateVehicle(vehicle *Vehicle) error {
	if vehicle == nil {
		return fmt.Errorf("vehicle 不能为 nil")
	}

	if vehicle.PlateNumber == "" {
		return fmt.Errorf("车牌号不能为空")
	}

	var routeName string
	if vehicle.RouteID != "" {
		route, err := GetRouteById(vehicle.RouteID)
		if err != nil {
			return err
		}
		if route == nil {
			return fmt.Errorf("线路不存在")
		}
		routeName = route.Name
	}

	now := util.GetMongoTimeNow()

	filter := bson.M{"plateNumber": vehicle.PlateNumber}
	update := bson.M{
		"$set": bson.M{
			"type":         vehicle.Type,
			"loadCapacity": vehicle.LoadCapacity,
			"status":       vehicle.Status,
			"routeId":      vehicle.RouteID,
			"routeName":    routeName,
			"remarks":      vehicle.Remarks,
			"lng":          vehicle.Lng,
			"lat":          vehicle.Lat,
			"currentLoad":  vehicle.CurrentLoad,
			"updateTime":   now,
		},
	}

	result, err := VehicleCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("未找到匹配的车辆进行更新")
	}

	return nil
}

// DeleteVehicle 删除车辆
func DeleteVehicle(plateNumber string) error {
	filter := bson.M{"plateNumber": plateNumber}
	_, err := VehicleCollection.DeleteOne(context.Background(), filter)
	return err
}

// GetVehicleList 根据条件查询车辆列表
func GetVehicleList(dto FindVehicleListDTO) (vehicles []*Vehicle, err error) {
	filter := bson.M{}
	if dto.PlateNumber != "" {
		filter["plateNumber"] = bson.M{"$regex": dto.PlateNumber, "$options": "i"}
	}
	if dto.Type != 0 {
		filter["type"] = dto.Type
	}
	if dto.Status != 0 {
		filter["status"] = dto.Status
	}
	if dto.RouteID != "" {
		filter["routeId"] = bson.M{"$regex": dto.RouteID, "$options": "i"}
	}
	if dto.RouteName != "" {
		filter["routeName"] = bson.M{"$regex": dto.RouteName, "$options": "i"}
	}
	findOptions := options.Find()
	findOptions.SetSkip(int64((dto.Page.Skip - 1) * dto.Page.Limit))
	findOptions.SetLimit(int64(dto.Page.Limit))
	findOptions.SetSort(bson.M{"updateTime": -1})

	cursor, err := VehicleCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var vehicle Vehicle
		if err := cursor.Decode(&vehicle); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, &vehicle)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return vehicles, nil
}

// GetVehicleTotalCount 获取车辆总数
func GetVehicleTotalCount(dto FindVehicleListDTO) (count int64, err error) {
	filter := bson.M{}
	if dto.PlateNumber != "" {
		filter["plateNumber"] = bson.M{"$regex": dto.PlateNumber, "$options": "i"}
	}
	if dto.Type != 0 {
		filter["type"] = dto.Type
	}
	if dto.Status != 0 {
		filter["status"] = dto.Status
	}
	if dto.RouteID != "" {
		filter["routeId"] = bson.M{"$regex": dto.RouteID, "$options": "i"}
	}
	if dto.RouteName != "" {
		filter["routeName"] = bson.M{"$regex": dto.RouteName, "$options": "i"}
	}
	documents, err := VehicleCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return
	}
	return documents, nil
}

// GetVehicleByRouteId 根据线路ID获取车辆列表（空闲车辆）
func GetVehicleByRouteId(routeId string) (vehicles []*Vehicle, err error) {
	filter := bson.M{}
	filter["routeId"] = routeId
	filter["status"] = Free
	cursor, err := VehicleCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var vehicle Vehicle
		if err := cursor.Decode(&vehicle); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, &vehicle)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return vehicles, nil
}

func GetVehicleById(plateNumber string) (vehicle *Vehicle, err error) {
	filter := bson.M{"plateNumber": plateNumber}
	err = VehicleCollection.FindOne(context.Background(), filter).Decode(&vehicle)
	return
}
