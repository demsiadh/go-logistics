package common

import "fmt"

// Page 分页参数结构体
type Page struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

func (p *Page) String() string {
	return fmt.Sprintf("skip: %d, limit: %d", p.Skip, p.Limit)
}

// GeoPoint 表示一个使用 GeoJSON 格式的坐标点。
type GeoPoint struct {
	Type        string    `bson:"type"`        // 类型，固定为 "Point"
	Coordinates []float64 `bson:"coordinates"` // 坐标，包含经度和纬度
}
