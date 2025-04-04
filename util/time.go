package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func GetMongoTimeNow() primitive.DateTime {
	currentTime := time.Now().UTC()
	return primitive.NewDateTimeFromTime(currentTime)
}
