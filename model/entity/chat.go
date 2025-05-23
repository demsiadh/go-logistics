package entity

import (
	"context"
	"encoding/json"
	"github.com/tmc/langchaingo/llms"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_logistics/config"
	"go_logistics/util"
)

var ChatCollection = config.MongoClient.Database("logistics").Collection("chat")

type Chat struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username   string             `bson:"username"`
	Title      string             `bson:"title"`
	Message    json.RawMessage    `bson:"message"`
	CreateTime primitive.DateTime `bson:"createTime" json:"createTime"`
	UpdateTime primitive.DateTime `bson:"updateTime" json:"updateTime"`
}

// ChatService service层与数据层分离
type ChatService struct {
	ID       string
	Username string
	Message  []llms.MessageContent
	Title    string
}

func GetChatListByUserName(username string) (chatServices []*ChatService, err error) {
	ctx := context.Background()
	filter := bson.M{"username": username}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"updateTime": -1})

	cursor, err := ChatCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var chat Chat
		if err = cursor.Decode(&chat); err != nil {
			return
		}

		chatServices = append(chatServices, &ChatService{
			ID:    chat.ID.Hex(),
			Title: chat.Title,
		})
	}

	if err = cursor.Err(); err != nil {
		return
	}
	return
}

func GetChatByIdAndUsername(id string, username string) (chatService *ChatService, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	filter := bson.M{
		"_id":      objectId,
		"username": username,
	}

	var chat Chat
	err = ChatCollection.FindOne(context.Background(), filter).Decode(&chat)
	if err != nil {
		return
	}

	// 反序列化 Message 字段为 []llms.MessageContent
	var messages []llms.MessageContent
	if chat.Message != nil {
		if err = json.Unmarshal(chat.Message, &messages); err != nil {
			return nil, err
		}
	}

	chatService = &ChatService{
		ID:       id,
		Username: chat.Username,
		Message:  messages,
		Title:    chat.Title,
	}

	return chatService, nil
}

func UpdateChatMessage(chatService *ChatService) (err error) {
	objectId, err := primitive.ObjectIDFromHex(chatService.ID)
	filter := bson.M{
		"_id":      objectId,
		"username": chatService.Username,
	}
	rawMessage, err := json.Marshal(chatService.Message)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"message":    rawMessage,
			"updateTime": util.GetMongoTimeNow(),
		},
	}
	_, err = ChatCollection.UpdateOne(context.Background(), filter, update)
	return
}

func InsertChat(chatService *ChatService) (err error) {
	now := util.GetMongoTimeNow()

	// 序列化 Message
	rawMessage, err := json.Marshal(chatService.Message)
	if err != nil {
		return err
	}

	chat := &Chat{
		ID:         primitive.NewObjectID(),
		Username:   chatService.Username,
		Title:      chatService.Title,
		Message:    rawMessage,
		CreateTime: now,
		UpdateTime: now,
	}

	_, err = ChatCollection.InsertOne(context.Background(), chat)
	return
}

func DeleteChat(chatId string) (err error) {
	objectId, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		return
	}
	_, err = ChatCollection.DeleteOne(context.Background(), bson.M{"_id": objectId})
	return
}

func UpdateChatTitle(chatId, username, title string) (err error) {
	objectId, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		return
	}
	filter := bson.M{
		"_id":      objectId,
		"username": username,
	}
	update := bson.M{
		"$set": bson.M{
			"title":      title,
			"updateTime": util.GetMongoTimeNow(),
		},
	}
	_, err = ChatCollection.UpdateOne(context.Background(), filter, update)
	return
}
