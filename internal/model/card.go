package model

import (
	"DonTaskMe-backend/internal/service"
	"context"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
)

type CardReq struct {
	Title   string `json:"title"`
	Index   int    `json:"index"`
	ListUID string `json:"listUid"`
}

type Card struct {
	Title string `json:"title"`
	Index int    `json:"index"`
	UID   string `json:"uid"`
}

func (c *CardReq) Save(ctx context.Context, listUID string) (*Card, error) {
	UID, _ := nano.Nanoid()
	newCard := Card{
		UID:   UID,
		Index: c.Index,
		Title: c.Title,
	}

	lh := service.DB.Collection(service.ListCollectionName)
	_, err := lh.UpdateOne(ctx, bson.M{"uid": listUID}, bson.D{{"$push", bson.D{{"cards", newCard}}}})
	if err != nil {
		return nil, err
	}
	return &newCard, err
}

func FindListCards(c context.Context, listUID string) ([]Card, error) {
	wh := service.DB.Collection(service.ListCollectionName)
	var list List
	err := wh.FindOne(c, bson.D{{"uid", listUID}}).Decode(&list)
	if err != nil {
		return nil, err
	}

	return list.Cards, nil
}

func DeleteCard(c context.Context, listUID string, cardUID string) error {
	wh := service.DB.Collection(service.WorkspaceCollectionName)
	res, err := wh.UpdateOne(
		c, bson.D{{"uid", listUID}},
		bson.D{{"$pull", bson.D{{"cards", bson.D{{"uid", cardUID}}}}}},
	)

	if err != nil {
		return err
	} else if res.ModifiedCount == 0 {
		return ResourceNotFound
	}
	return nil
}
