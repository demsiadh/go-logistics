package vo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go_logistics/model/entity"
)

type OrderVO struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	OrderID          string             `bson:"orderId" json:"orderId"`
	CustomerName     string             `bson:"customerName" json:"customerName"`
	Phone            string             `bson:"phone" json:"phone"`
	StartAddress     string             `bson:"startAddress" json:"startAddress"`
	StartLng         string             `bson:"startLng" json:"startLng"`
	StartLat         string             `bson:"startLat" json:"startLat"`
	StartOutlet      entity.Outlet      `bson:"startOutlet" json:"startOutlet"`
	EndAddress       string             `bson:"endAddress" json:"endAddress"`
	EndLng           string             `bson:"endLng" json:"endLng"`
	EndLat           string             `bson:"endLat" json:"endLat"`
	EndOutlet        entity.Outlet      `bson:"endOutlet" json:"endOutlet"`
	TransPortVehicle entity.Vehicle     `bson:"transPortVehicle" json:"transPortVehicle"`
	Route            entity.Route       `bson:"route" json:"route"`
	Weight           float64            `bson:"weight" json:"weight"`
	Status           entity.OrderStatus `bson:"status" json:"status"`
	CreateTime       primitive.DateTime `bson:"createTime" json:"createTime"`
	UpdateTime       primitive.DateTime `bson:"updateTime" json:"-"`
	Remark           string             `bson:"remark" json:"remark"`
}

func ToOrderVO(order *entity.Order) (OrderVO, error) {
	var startOutlet *entity.Outlet
	var endOutlet *entity.Outlet
	var vehicle *entity.Vehicle
	var route *entity.Route

	// 处理开始网点
	if order.StartOutletId != "" {
		outlet, err := entity.GetOutletById(order.StartOutletId)
		if err != nil {
			// 可选：记录日志
			// log.Printf("获取开始网点失败: %v", err)
		} else {
			startOutlet = outlet
		}
	}

	// 处理结束网点
	if order.EndOutletId != "" {
		outlet, err := entity.GetOutletById(order.EndOutletId)
		if err != nil {
			// log.Printf("获取结束网点失败: %v", err)
		} else {
			endOutlet = outlet
		}
	}

	// 处理运输车辆
	if order.TransPortVehicle != "" {
		v, err := entity.GetVehicleById(order.TransPortVehicle)
		if err != nil {
			// log.Printf("获取车辆失败: %v", err)
		} else {
			vehicle = v
		}
	}

	// 处理线路（依赖车辆）
	if vehicle != nil && vehicle.RouteID != "" {
		r, err := entity.GetRouteById(vehicle.RouteID)
		if err != nil {
			// log.Printf("获取线路失败: %v", err)
		} else {
			route = r
		}
	}

	orderVO := OrderVO{
		ID:           order.ID,
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		Phone:        order.Phone,
		StartAddress: order.StartAddress,
		StartLng:     order.StartLng,
		StartLat:     order.StartLat,
		StartOutlet: func() entity.Outlet {
			if startOutlet != nil {
				return *startOutlet
			}
			return entity.Outlet{}
		}(),
		EndAddress: order.EndAddress,
		EndLng:     order.EndLng,
		EndLat:     order.EndLat,
		EndOutlet: func() entity.Outlet {
			if endOutlet != nil {
				return *endOutlet
			}
			return entity.Outlet{}
		}(),
		TransPortVehicle: func() entity.Vehicle {
			if vehicle != nil {
				return *vehicle
			}
			return entity.Vehicle{}
		}(),
		Route: func() entity.Route {
			if route != nil {
				return *route
			}
			return entity.Route{}
		}(),
		Weight:     order.Weight,
		Status:     order.Status,
		CreateTime: order.CreateTime,
		UpdateTime: order.UpdateTime,
		Remark:     order.Remark,
	}
	return orderVO, nil
}

func ToOrderVOList(orders []*entity.Order) ([]OrderVO, error) {
	orderVOs := make([]OrderVO, 0)
	for _, order := range orders {
		orderVO, err := ToOrderVO(order)
		if err != nil {
			return nil, err
		}
		orderVOs = append(orderVOs, orderVO)
	}
	return orderVOs, nil
}
