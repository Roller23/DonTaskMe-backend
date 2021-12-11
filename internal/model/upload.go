package model

import (
	"DonTaskMe-backend/internal/service"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type FileInfo struct {
	Filename    string `json:"filename"`
	StoragePath string `json:"storagePath"`
}

func (f *FileInfo) Save(c context.Context, cardUID string) error {
	lh := service.DB.Collection(service.ListCollectionName)
	filter := bson.M{"cards": bson.M{"$elemMatch": bson.M{"uid": cardUID}}}
	fileInfo := bson.D{{"cards.$.files", f}}
	update := bson.D{{"$push", fileInfo}}
	res, err := lh.UpdateOne(c, filter, update)
	if err != nil {
		return err
	} else if res.ModifiedCount == 0 {
		return ResourceNotFound
	}
	return nil
}
