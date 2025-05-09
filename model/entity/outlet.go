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

var OutletCollection = config.MongoClient.Database("logistics").Collection("outlet")

// OutletStatus 网点状态的枚举
type OutletStatus int

const (
	OutletStatusOpen   OutletStatus = 1 // 营业中
	OutletStatusClosed OutletStatus = 2 // 已关闭
)

func (s OutletStatus) String() string {
	return [...]string{"营业中", "已关闭"}[s-1]
}

type Outlet struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Phone         string             `bson:"phone" json:"phone"`
	Province      string             `bson:"province" json:"province"`
	City          string             `bson:"city" json:"city"`
	DetailAddress string             `bson:"detailAddress" json:"detailAddress"`
	BusinessHours string             `bson:"businessHours" json:"businessHours"`
	Lng           string             `bson:"lng" json:"lng"`
	Lat           string             `bson:"lat" json:"lat"`
	Scope         []common.GeoPoint  `bson:"scope" json:"scope"`
	Status        OutletStatus       `bson:"status" json:"status"`
	Remark        string             `bson:"remark" json:"remark"`
	CreateTime    primitive.DateTime `bson:"createTime" json:"-"`
	UpdateTime    primitive.DateTime `bson:"updateTime" json:"-"`
}

// FindOutletListDTO 查询网点列表的参数
type FindOutletListDTO struct {
	Name     string       `json:"name"`
	Status   OutletStatus `json:"status"`
	Province string       `json:"province"`
	City     string       `json:"city"`
	Page     common.Page  `json:"page"`
}

func (dto *FindOutletListDTO) String() string {
	return fmt.Sprintf("name: %s, status: %s, page: %s", dto.Name, dto.Status, dto.Page.String())
}

// InsertOutlet 新建网点
func InsertOutlet(outlet *Outlet) error {
	// 填充时间
	outlet.CreateTime = util.GetMongoTimeNow()
	outlet.UpdateTime = util.GetMongoTimeNow()
	_, err := OutletCollection.InsertOne(context.Background(), outlet)
	return err
}

// UpdateOutlet 修改网点信息
func UpdateOutlet(outletId string, outlet *Outlet) error {
	// 将 outletId 转换为 primitive.ObjectID
	objectId, err := primitive.ObjectIDFromHex(outletId)
	if err != nil {
		return fmt.Errorf("invalid outletId: %w", err)
	}

	// 构建过滤条件
	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": bson.M{
			"name":          outlet.Name,
			"phone":         outlet.Phone,
			"detailAddress": outlet.DetailAddress,
			"businessHours": outlet.BusinessHours,
			"lng":           outlet.Lng,
			"lat":           outlet.Lat,
			"scope":         outlet.Scope,
			"status":        outlet.Status,
			"remark":        outlet.Remark,
			"updateTime":    util.GetMongoTimeNow(),
		},
	}
	_, err = OutletCollection.UpdateOne(context.Background(), filter, update)
	return err
}

// DeleteOutlet 删除网点
func DeleteOutlet(outletId string) error {
	// 将 outletId 转换为 primitive.ObjectID
	objectId, err := primitive.ObjectIDFromHex(outletId)
	if err != nil {
		return fmt.Errorf("invalid outletId: %w", err)
	}

	// 构建过滤条件
	filter := bson.M{"_id": objectId}

	// 执行删除操作
	_, err = OutletCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

// GetOutletList 根据条件查询网点列表
func GetOutletList(dto FindOutletListDTO) (outlets []*Outlet, err error) {
	filter := bson.M{}
	if dto.Name != "" {
		filter["name"] = bson.M{"$regex": dto.Name, "$options": "i"}
	}
	if dto.Province != "" {
		filter["province"] = bson.M{"$regex": dto.Province, "$options": "i"}
	}
	if dto.City != "" {
		filter["city"] = bson.M{"$regex": dto.City, "$options": "i"}
	}
	if dto.Status != 0 {
		filter["status"] = dto.Status
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((dto.Page.Skip - 1) * dto.Page.Limit))
	findOptions.SetLimit(int64(dto.Page.Limit))
	findOptions.SetSort(bson.M{"name": 1})

	cursor, err := OutletCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var Outlet Outlet
		if err := cursor.Decode(&Outlet); err != nil {
			return nil, err
		}
		outlets = append(outlets, &Outlet)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return outlets, nil
}

// GetOutletTotalCount 获取总数
func GetOutletTotalCount() (count int64, err error) {
	documents, err := OutletCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return
	}
	return documents, nil
}

// GetAllProvincesAndCities 查询数据库中所有的省份和城市列表
func GetAllProvincesAndCities() (result map[string][]string, err error) {
	result = make(map[string][]string)

	// 查询所有不重复的省份
	provincePipeline := []bson.M{
		{"$group": bson.M{"_id": "$province"}},
		{"$project": bson.M{"_id": 0, "province": "$_id"}},
	}
	provinceCursor, err := OutletCollection.Aggregate(context.Background(), provincePipeline)
	if err != nil {
		return nil, err
	}
	defer provinceCursor.Close(context.Background())
	for provinceCursor.Next(context.Background()) {
		var provinceDoc struct {
			Province string `bson:"province"`
		}
		if err := provinceCursor.Decode(&provinceDoc); err != nil {
			return nil, err
		}
		result["provinces"] = append(result["provinces"], provinceDoc.Province)
	}
	if err := provinceCursor.Err(); err != nil {
		return nil, err
	}

	// 查询所有不重复的城市
	cityPipeline := []bson.M{
		{"$group": bson.M{"_id": "$city"}},
		{"$project": bson.M{"_id": 0, "city": "$_id"}},
	}
	cityCursor, err := OutletCollection.Aggregate(context.Background(), cityPipeline)
	if err != nil {
		return nil, err
	}
	defer cityCursor.Close(context.Background())
	for cityCursor.Next(context.Background()) {
		var cityDoc struct {
			City string `bson:"city"`
		}
		if err := cityCursor.Decode(&cityDoc); err != nil {
			return nil, err
		}
		result["cities"] = append(result["cities"], cityDoc.City)
	}
	if err := cityCursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func GetOutletById(outletId string) (outlet *Outlet, err error) {
	var objectId, _ = primitive.ObjectIDFromHex(outletId)
	filter := bson.M{"_id": objectId}
	err = OutletCollection.FindOne(context.Background(), filter).Decode(&outlet)
	return
}
