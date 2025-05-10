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
	startOutlet, err := entity.GetOutletById(order.StartOutletId)
	if err != nil {
		return OrderVO{}, err
	}
	endOutlet, err := entity.GetOutletById(order.EndOutletId)
	if err != nil {
		return OrderVO{}, err
	}
	vehicle, err := entity.GetVehicleById(order.TransPortVehicle)
	if err != nil {
		return OrderVO{}, err
	}
	route, err := entity.GetRouteById(vehicle.RouteID)
	if err != nil {
		return OrderVO{}, err
	}
	return OrderVO{
		ID:               order.ID,
		OrderID:          order.OrderID,
		CustomerName:     order.CustomerName,
		Phone:            order.Phone,
		StartAddress:     order.StartAddress,
		StartLng:         order.StartLng,
		StartLat:         order.StartLat,
		StartOutlet:      *startOutlet,
		EndAddress:       order.EndAddress,
		EndLng:           order.EndLng,
		EndLat:           order.EndLat,
		EndOutlet:        *endOutlet,
		TransPortVehicle: *vehicle,
		Route:            *route,
		Weight:           order.Weight,
		Status:           order.Status,
		CreateTime:       order.CreateTime,
		UpdateTime:       order.UpdateTime,
		Remark:           order.Remark,
	}, nil
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
