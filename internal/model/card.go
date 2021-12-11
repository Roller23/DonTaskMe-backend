package model

import (
	"DonTaskMe-backend/internal/service"
	"context"

	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
)

type CardReq struct {
	Title       string `json:"title"`
	Index       int    `json:"index"`
	ListUID     string `json:"listUid"`
	Description string `json:"description"`
}

type CardUpdateReq struct {
	Title       *string `json:"title,omitempty"`
	Index       *int    `json:"index,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Card struct {
	Title       string `json:"title"`
	Index       int    `json:"index"`
	UID         string `json:"uid"`
	Description string `json:"description"`
}

func (c *CardReq) Save(ctx context.Context, listUID string) (*Card, error) {
	UID, _ := nano.Nanoid()
	newCard := Card{
		UID:         UID,
		Index:       c.Index,
		Title:       c.Title,
		Description: c.Description,
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
	wh := service.DB.Collection(service.ListCollectionName)
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

func UpdateCard(c context.Context, cardUID string, updateBody CardUpdateReq) error {
	wh := service.DB.Collection(service.ListCollectionName)
	//{awards: {$elemMatch: {award:'National Medal', year:1975}}}
	filter := bson.M{"cards": bson.M{"$elemMatch": bson.M{"uid": cardUID}}}
	var toSet bson.D
	combineIfExists(&toSet, "index", updateBody.Index)
	combineIfExists(&toSet, "title", updateBody.Title)
	combineIfExists(&toSet, "description", updateBody.Description)
	update := bson.D{{"$set", toSet}}
	res, err := wh.UpdateOne(c, filter, update)
	if err != nil {
		return err
	} else if res.ModifiedCount == 0 {
		return ResourceNotFound
	}
	return nil
}

func combineIfExists(doc *bson.D, key string, val interface{}) {
	if val != nil {
		*doc = append(*doc, bson.E{Key: key, Value: val})
	}
}
