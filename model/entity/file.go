package entity

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_logistics/common"
	"go_logistics/config"
	"go_logistics/util"
)

var FileCollection = config.MongoClient.Database("logistics").Collection("file")

// BusinessType 文件所属业务类型
type BusinessType int

const (
	AIRepository BusinessType = 1
)

func (bt BusinessType) String() string {
	textMap := map[BusinessType]string{
		AIRepository: "AI知识库",
	}
	return textMap[bt]
}

type File struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileType    BusinessType       `bson:"fileType" json:"fileType"`
	FileName    string             `bson:"fileName" json:"fileName"`
	FileSize    int64              `bson:"fileSize" json:"fileSize"`
	ContentType string             `bson:"contentType" json:"contentType"`
	FileData    []byte             `bson:"fileData" json:"fileData"`
	VectorIds   []string           `bson:"vectorIds" json:"VectorIds"`
	UploadTime  primitive.DateTime `bson:"uploadTime" json:"uploadTime"`
}

type FindFileListDTO struct {
	FileType BusinessType `json:"fileType"`
	FileName string       `json:"fileName"`
	Page     common.Page  `json:"page"`
}

func InsertFile(ctx context.Context, file *File) (err error) {
	file.UploadTime = util.GetMongoTimeNow()
	_, err = FileCollection.InsertOne(ctx, file)
	return err
}

func DeleteFile(ctx context.Context, fileId string) (err error) {
	objectId, err := primitive.ObjectIDFromHex(fileId)
	if err != nil {
		return
	}
	_, _ = FileCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	return
}

func GetFileList(ctx context.Context, dto FindFileListDTO) (files []*File, err error) {
	filter := bson.M{}
	if dto.FileType != 0 {
		filter["fileType"] = dto.FileType
	}
	if dto.FileName != "" {
		filter["fileName"] = bson.M{"$regex": dto.FileName, "$options": "i"}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"uploadTime": -1})
	findOptions.SetSkip(int64((dto.Page.Skip - 1) * dto.Page.Limit))
	findOptions.SetLimit(int64(dto.Page.Limit))

	cursor, err := FileCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var file File
		if err = cursor.Decode(&file); err != nil {
			return
		}
		files = append(files, &file)
	}
	if err = cursor.Err(); err != nil {
		return
	}
	return
}

func GetFileById(ctx context.Context, fileId string) (file *File, err error) {
	objectId, err := primitive.ObjectIDFromHex(fileId)
	if err != nil {
		return
	}
	err = FileCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&file)
	return
}
