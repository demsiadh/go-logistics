package vo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go_logistics/model/entity"
)

type VehicleVO struct {
	ID           string               `bson:"_id,omitempty" json:"id"`
	PlateNumber  string               `bson:"plateNumber" json:"plateNumber"`
	Type         entity.VehicleType   `bson:"type" json:"type"`
	LoadCapacity float64              `bson:"loadCapacity" json:"loadCapacity"`
	CurrentLoad  float64              `bson:"currentLoad" json:"currentLoad"`
	Status       entity.VehicleStatus `bson:"status" json:"status"`
	RouteID      string               `bson:"routeId" json:"routeId"`
	RouteName    string               `bson:"routeName" json:"routeName"`
	Remarks      string               `bson:"remarks" json:"remarks"`
	Lng          string               `bson:"lng" json:"lng"`
	Lat          string               `bson:"lat" json:"lat"`
	Route        *entity.Route        `bson:"route" json:"route"`
	CreateTime   primitive.DateTime   `bson:"createTime" json:"-"`
	UpdateTime   primitive.DateTime   `bson:"updateTime" json:"-"`
}

func ToVehicleVO(vehicle *entity.Vehicle) (VehicleVO, error) {
	var route *entity.Route

	if vehicle.RouteID != "" {
		var err error
		route, err = entity.GetRouteById(vehicle.RouteID)
		if err != nil {
			return VehicleVO{}, err
		}
	}

	return VehicleVO{
		ID:           vehicle.ID,
		PlateNumber:  vehicle.PlateNumber,
		Type:         vehicle.Type,
		LoadCapacity: vehicle.LoadCapacity,
		CurrentLoad:  vehicle.CurrentLoad,
		Status:       vehicle.Status,
		RouteID:      vehicle.RouteID,
		RouteName:    vehicle.RouteName,
		Remarks:      vehicle.Remarks,
		Lng:          vehicle.Lng,
		Lat:          vehicle.Lat,
		Route:        route,
		CreateTime:   vehicle.CreateTime,
		UpdateTime:   vehicle.UpdateTime,
	}, nil
}

func ToVehicleVOList(vehicles []*entity.Vehicle) ([]VehicleVO, error) {
	var vehicleVOs []VehicleVO

	for _, vehicle := range vehicles {
		vehicleVO, err := ToVehicleVO(vehicle)
		if err != nil {
			return nil, err
		}
		vehicleVOs = append(vehicleVOs, vehicleVO)
	}

	return vehicleVOs, nil
}
