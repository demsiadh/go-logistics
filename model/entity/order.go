package entity

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_logistics/common"
	"go_logistics/config"
	"go_logistics/util"
	"time"
)

var OrderCollection = config.MongoClient.Database("logistics").Collection("order")

// OrderStatus 订单状态
type OrderStatus int

const (
	Pending    OrderStatus = 1
	Processing OrderStatus = 2
	Completed  OrderStatus = 3
	Cancelled  OrderStatus = 4
)

func (s OrderStatus) String() string {
	textMap := map[OrderStatus]string{
		Pending:    "待处理",
		Processing: "处理中",
		Completed:  "已完成",
		Cancelled:  "已取消",
	}
	return textMap[s]
}

// Order 订单结构
type Order struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	OrderID          string             `bson:"orderId" json:"orderId"`
	CustomerName     string             `bson:"customerName" json:"customerName"`
	Phone            string             `bson:"phone" json:"phone"`
	StartAddress     string             `bson:"startAddress" json:"startAddress"`
	StartLng         string             `bson:"startLng" json:"startLng"`
	StartLat         string             `bson:"startLat" json:"startLat"`
	StartOutletId    string             `bson:"startOutletId" json:"startOutletId"`
	EndAddress       string             `bson:"endAddress" json:"endAddress"`
	EndLng           string             `bson:"endLng" json:"endLng"`
	EndLat           string             `bson:"endLat" json:"endLat"`
	EndOutletId      string             `bson:"endOutletId" json:"endOutletId"`
	TransPortVehicle string             `bson:"transPortVehicle" json:"transPortVehicle"`
	Weight           float64            `bson:"weight" json:"weight"`
	Status           OrderStatus        `bson:"status" json:"status"`
	CreateTime       primitive.DateTime `bson:"createTime" json:"createTime"`
	UpdateTime       primitive.DateTime `bson:"updateTime" json:"-"`
	Remark           string             `bson:"remark" json:"remark"`
}

// FindOrderListDTO 查询订单列表的参数
type FindOrderListDTO struct {
	OrderID   string      `json:"orderId"`
	Phone     string      `json:"phone"`
	Status    OrderStatus `json:"status"`
	StartTime time.Time   `json:"startTime"`
	EndTime   time.Time   `json:"endTime"`
	Page      common.Page `json:"page"`
}

func (dto *FindOrderListDTO) String() string {
	return fmt.Sprintf("orderId: %s, phone: %s, status: %d, startTime: %s, endTime: %s, page: %s",
		dto.OrderID, dto.Phone, dto.Status, dto.StartTime, dto.EndTime, dto.Page.String())
}

// InsertOrder 新建订单
func InsertOrder(order *Order) error {
	// 填充时间
	order.CreateTime = util.GetMongoTimeNow()
	order.UpdateTime = util.GetMongoTimeNow()
	_, err := OrderCollection.InsertOne(context.Background(), order)
	return err
}

// UpdateOrder 修改订单信息
func UpdateOrder(order *Order) error {
	filter := bson.M{"orderId": order.OrderID}
	update := bson.M{
		"$set": bson.M{
			"customerName": order.CustomerName,
			"phone":        order.Phone,
			"status":       order.Status,
			"remark":       order.Remark,
			"updateTime":   util.GetMongoTimeNow(),
		},
	}
	_, err := OrderCollection.UpdateOne(context.Background(), filter, update)
	return err
}

// CompleteDataOrder 完成订单数据
func CompleteDataOrder(order *Order) error {
	filter := bson.M{"orderId": order.OrderID}
	update := bson.M{
		"$set": bson.M{
			"status":           Processing,
			"startOutletId":    order.StartOutletId,
			"endOutletId":      order.EndOutletId,
			"transPortVehicle": order.TransPortVehicle,
			"updateTime":       util.GetMongoTimeNow(),
			"remark":           order.Remark,
		},
	}
	_, err := OrderCollection.UpdateOne(context.Background(), filter, update)
	return err
}

// DeleteOrder 删除订单
func DeleteOrder(orderId string) error {
	filter := bson.M{"orderId": orderId}
	_, err := OrderCollection.DeleteOne(context.Background(), filter)
	return err
}

// GetOrderList 根据条件查询订单列表
func GetOrderList(dto FindOrderListDTO) (orders []*Order, err error) {
	filter := bson.M{}
	if dto.OrderID != "" {
		filter["orderId"] = bson.M{"$regex": dto.OrderID, "$options": "i"}
	}

	if dto.Phone != "" {
		filter["phone"] = bson.M{"$regex": dto.Phone, "$options": "i"}
	}
	if dto.Status != 0 {
		filter["status"] = dto.Status
	}
	if !dto.StartTime.IsZero() || !dto.EndTime.IsZero() {
		timeFilter := bson.M{}
		if !dto.StartTime.IsZero() {
			// 将本地时间转为 UTC 时间
			startUTC := dto.StartTime.UTC()
			timeFilter["$gte"] = primitive.NewDateTimeFromTime(startUTC)
		}
		if !dto.EndTime.IsZero() {
			// 将本地时间转为 UTC 时间
			endUTC := dto.EndTime.UTC()
			timeFilter["$lte"] = primitive.NewDateTimeFromTime(endUTC)
		}
		filter["createTime"] = timeFilter
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((dto.Page.Skip - 1) * dto.Page.Limit))
	findOptions.SetLimit(int64(dto.Page.Limit))
	findOptions.SetSort(bson.M{"updateTime": -1})

	cursor, err := OrderCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var order Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderTotalCount 获取订单总数
func GetOrderTotalCount(dto FindOrderListDTO) (count int64, err error) {
	filter := bson.M{}
	if dto.OrderID != "" {
		filter["orderId"] = bson.M{"$regex": dto.OrderID, "$options": "i"}
	}

	if dto.Phone != "" {
		filter["phone"] = bson.M{"$regex": dto.Phone, "$options": "i"}
	}
	if dto.Status != 0 {
		filter["status"] = dto.Status
	}
	if !dto.StartTime.IsZero() || !dto.EndTime.IsZero() {
		timeFilter := bson.M{}
		if !dto.StartTime.IsZero() {
			timeFilter["$gte"] = primitive.NewDateTimeFromTime(dto.StartTime)
		}
		if !dto.EndTime.IsZero() {
			timeFilter["$lte"] = primitive.NewDateTimeFromTime(dto.EndTime)
		}
		filter["createTime"] = timeFilter
	}
	documents, err := OrderCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return
	}
	return documents, nil
}

// GetOrderCountByDate 获取指定日期的订单数量（适配中国时区）
func GetOrderCountByDate(date string) (int, error) {
	// 解析输入日期，格式如 "20060102"
	parsedDate, err := time.ParseInLocation("20060102", date, time.Local)
	if err != nil {
		return 0, err
	}

	// 构造当天的起止时间（本地时间，即中国时间）
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	// 转换为 UTC 时间用于数据库查询
	startUTC := startOfDay.UTC()
	endUTC := endOfDay.UTC()

	filter := bson.M{
		"createTime": bson.M{
			"$gte": primitive.NewDateTimeFromTime(startUTC),
			"$lt":  primitive.NewDateTimeFromTime(endUTC),
		},
	}

	count, err := OrderCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func GetOrderById(orderId string) (*Order, error) {
	filter := bson.M{"orderId": orderId}
	var order Order
	err := OrderCollection.FindOne(context.Background(), filter).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CompleteOrderByVehicle 以车辆为单位完成订单
func CompleteOrderByVehicle(vehicle *Vehicle) error {
	filter := bson.M{
		"transPortVehicle": vehicle.PlateNumber,
		"status":           Processing,
	}
	update := bson.M{
		"$set": bson.M{
			"status":     Completed,
			"updateTime": util.GetMongoTimeNow(),
		},
	}
	_, err := OrderCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// GetOrderByVehicle 根据车牌号获取订单
func GetOrderByVehicle(plateNumber string) (order *Order, err error) {
	filter := bson.M{
		"transPortVehicle": plateNumber,
		"status":           Processing,
	}

	var result Order
	err = OrderCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 可选：将“未找到”视为无错误，返回 nil, nil
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
