package model

import (
	"DonTaskMe-backend/internal/service"
	"context"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type CommentReq struct {
	Content string `json:"content"`
}

type Comment struct {
	UID      string `json:"uid"`
	Content  string `json:"content"`
	Date     int64  `json:"date"`
	Username string `json:"username"`
}

func (req *CommentReq) Save(c context.Context, cardUID, username string) (*Comment, error) {
	UID, _ := nano.Nanoid()
	newComment := Comment{
		UID:      UID,
		Content:  req.Content,
		Date:     time.Now().Unix(),
		Username: username,
	}

	lh := service.DB.Collection(service.ListCollectionName)
	filter := bson.M{"cards": bson.M{"$elemMatch": bson.M{"uid": cardUID}}}
	commentDoc := bson.D{{"cards.$.comments", newComment}}
	update := bson.D{{"$push", commentDoc}}
	res, err := lh.UpdateOne(c, filter, update)
	if err != nil {
		return nil, err
	} else if res.ModifiedCount == 0 {
		return nil, ResourceNotFound
	}
	return &newComment, nil
}
