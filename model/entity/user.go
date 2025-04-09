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

var UserCollection = config.MongoClient.Database("logistics").Collection("user")

// UserStatus 用户状态
type UserStatus int

const (
	Active UserStatus = iota + 1
	Banned
	Deleted
)

func (s UserStatus) String() string {
	return [...]string{"活跃", "禁止", "删除"}[s-1]
}

// User 用户结构
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Name       string             `bson:"name" json:"name"`
	Phone      string             `bson:"phone" json:"phone"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"-"`
	Status     UserStatus         `bson:"status" json:"status"`
	Salt       string             `bson:"salt" json:"-"`
	CreateTime primitive.DateTime `bson:"createTime" json:"-"`
	UpdateTime primitive.DateTime `bson:"updateTime" json:"-"`
}

// FindUserListDTO 查询用户列表的参数
type FindUserListDTO struct {
	Name   string      `json:"name"`
	Phone  string      `json:"phone"`
	Email  string      `json:"email"`
	Status UserStatus  `json:"status"`
	Page   common.Page `json:"page"`
}

func (dto *FindUserListDTO) String() string {
	return fmt.Sprintf("name: %s, phone: %s, email: %s, status: %d, page: %s", dto.Name, dto.Phone, dto.Email, dto.Status, dto.Page.String())
}

// InsertUser 新建用户
func InsertUser(user *User) error {
	// 填充时间
	user.CreateTime = util.GetMongoTimeNow()
	user.UpdateTime = util.GetMongoTimeNow()
	_, err := UserCollection.InsertOne(context.Background(), user)
	return err
}

// UpdateUser 修改用户信息
func UpdateUser(user *User) error {
	user.UpdateTime = util.GetMongoTimeNow()
	filter := bson.M{"name": user.Name}
	update := bson.M{
		"$set": bson.M{
			"phone":      user.Phone,
			"email":      user.Email,
			"status":     user.Status,
			"updateTime": user.UpdateTime,
		},
	}
	_, err := UserCollection.UpdateOne(context.Background(), filter, update)
	return err
}

// DeleteUser 删除用户（逻辑删除）
func DeleteUser(name string) error {
	fmt.Println(util.GetMongoTimeNow().Time())
	update := bson.M{
		"$set": bson.M{
			"status":     Deleted,
			"updateTime": util.GetMongoTimeNow(),
		},
	}
	filter := bson.M{"name": name}
	_, err := UserCollection.UpdateOne(context.Background(), filter, update)
	return err
}

// GetUserByName 根据用户名查询用户信息
func GetUserByName(name string) (*User, error) {
	var user User
	filter := bson.M{"name": name}
	err := UserCollection.FindOne(context.Background(), filter).Decode(&user)
	return &user, err
}

// GetUserList 根据条件查询用户列表
func GetUserList(dto FindUserListDTO) (users []*User, err error) {
	filter := bson.M{}
	if dto.Name != "" {
		filter["name"] = bson.M{"$regex": "^" + dto.Name, "$options": "i"}
	}
	if dto.Phone != "" {
		filter["phone"] = bson.M{"$regex": "^" + dto.Phone, "$options": "i"}
	}
	if dto.Email != "" {
		filter["email"] = bson.M{"$regex": "^" + dto.Email, "$options": "i"}
	}
	if dto.Status != 0 {
		filter["status"] = dto.Status
	}
	findOptions := options.Find()
	findOptions.SetSkip(int64((dto.Page.Skip - 1) * dto.Page.Limit))
	findOptions.SetLimit(int64(dto.Page.Limit))
	findOptions.SetSort(bson.M{"name": 1})

	cursor, err := UserCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
