package util

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"go_logistics/common"
	"math"
	"strconv"
)

const earthRadiusKm = 6371.0

// haversin(θ) = sin²(θ/2)
func haversin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// GetDistanceFromString 接收经纬度字符串并计算球面距离（单位：公里）
func GetDistanceFromString(lat1Str, lng1Str, lat2Str, lng2Str string) (float64, error) {
	lat1, err := strconv.ParseFloat(lat1Str, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid lat1: %v", err)
	}
	lng1, err := strconv.ParseFloat(lng1Str, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid lng1: %v", err)
	}
	lat2, err := strconv.ParseFloat(lat2Str, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid lat2: %v", err)
	}
	lng2, err := strconv.ParseFloat(lng2Str, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid lng2: %v", err)
	}

	return GetDistance(lat1, lng1, lat2, lng2), nil
}

// GetDistance 计算两个经纬度点之间的球面距离（单位：公里）
func GetDistance(lat1, lng1, lat2, lng2 float64) float64 {
	rad := func(d float64) float64 { return d * math.Pi / 180 }

	lat1 = rad(lat1)
	lng1 = rad(lng1)
	lat2 = rad(lat2)
	lng2 = rad(lng2)

	dLat := lat2 - lat1
	dLng := lng2 - lng1

	a := haversin(dLat) + math.Cos(lat1)*math.Cos(lat2)*haversin(dLng)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// IsPointInGeoPointSlice 查询当前点位是否在范围中
func IsPointInGeoPointSlice(targetLngStr, targetLatStr string, points []common.GeoPoint) (bool, error) {
	targetLat, err := strconv.ParseFloat(targetLatStr, 64)
	if err != nil {
		return false, fmt.Errorf("invalid lat1: %v", err)
	}
	targetLng, err := strconv.ParseFloat(targetLngStr, 64)
	if err != nil {
		return false, fmt.Errorf("invalid lng1: %v", err)
	}
	if len(points) < 3 {
		return false, nil // 至少需要三个点才能构成面
	}

	// 构造 orb.Ring
	var ring orb.Ring
	for _, p := range points {
		// 假设 Coordinates 中顺序是 [lng, lat]
		ring = append(ring, orb.Point{p.Coordinates[0], p.Coordinates[1]})
	}

	poly := orb.Polygon{ring}
	point := orb.Point{targetLng, targetLat}

	return planar.PolygonContains(poly, point), nil
}
