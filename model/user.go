package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go_logistics/config"
)

var UserCollection = config.MongoClient.Database("logistics").Collection("user")

// UserStatus 用户状态
type UserStatus int

const (
	Active UserStatus = iota
	Inactive
	Banned
	Deleted
)

func (s UserStatus) String() string {
	return [...]string{"活跃", "不活跃", "禁止", "删除"}[s]
}

// User 用户结构
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Name     string             `bson:"name"`
	Phone    string             `bson:"phone"`
	Email    string             `bson:"email"`
	Password string             `bson:"password" json:"-"`
	Status   UserStatus         `bson:"status"`
	Salt     string             `bson:"salt" json:"-"`
}

func InsertUser(user *User) error {
	_, err := UserCollection.InsertOne(context.Background(), user)
	return err
}

func UpdateUser(user *User) error {
	_, err := UserCollection.UpdateOne(context.Background(), user.Name, user)
	return err
}

func DeleteUser(user *User) error {
	_, err := UserCollection.DeleteOne(context.Background(), user.Name)
	return err
}

func GetUserByName(name string) (*User, error) {
	var user User
	filter := bson.M{"name": name}
	err := UserCollection.FindOne(context.Background(), filter).Decode(&user)
	return &user, err
}
