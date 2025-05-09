package vo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go_logistics/common"
	"go_logistics/model/entity"
)

type RouteVO struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RouteID     string             `bson:"routeId" json:"routeId"`
	Name        string             `bson:"name" json:"name"`
	Type        entity.RouteType   `bson:"type" json:"type"`
	Status      entity.RouteStatus `bson:"status" json:"status"`
	Description string             `bson:"description" json:"description"`
	Points      []common.GeoPoint  `bson:"points" json:"points"`
	Distance    float64            `bson:"distance" json:"distance"`
	StartOutlet *entity.Outlet     `bson:"startOutlet" json:"startOutlet"`
	EndOutlet   *entity.Outlet     `bson:"endOutlet" json:"endOutlet"`
	CreateTime  primitive.DateTime `bson:"createTime" json:"-"`
	UpdateTime  primitive.DateTime `bson:"updateTime" json:"-"`
}

func ToRouteVO(route *entity.Route) (RouteVO, error) {
	var startOutlet *entity.Outlet
	var endOutlet *entity.Outlet

	if route.StartOutlet != "" {
		var err error
		startOutlet, err = entity.GetOutletById(route.StartOutlet)
		if err != nil {
			return RouteVO{}, err
		}
	}

	if route.EndOutlet != "" {
		var err error
		endOutlet, err = entity.GetOutletById(route.EndOutlet)
		if err != nil {
			return RouteVO{}, err
		}
	}

	return RouteVO{
		ID:          route.ID,
		RouteID:     route.RouteID,
		Name:        route.Name,
		Type:        route.Type,
		Status:      route.Status,
		Description: route.Description,
		Points:      route.Points,
		Distance:    route.Distance,
		StartOutlet: startOutlet,
		EndOutlet:   endOutlet,
		CreateTime:  route.CreateTime,
		UpdateTime:  route.UpdateTime,
	}, nil
}

func ToRouteVOList(routes []*entity.Route) ([]RouteVO, error) {
	var routeVOs []RouteVO
	for _, route := range routes {
		routeVO, err := ToRouteVO(route)
		if err != nil {
			return nil, err
		}
		routeVOs = append(routeVOs, routeVO)
	}
	return routeVOs, nil
}
