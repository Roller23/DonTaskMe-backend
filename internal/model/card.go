package model

import (
	"DonTaskMe-backend/internal/service"
	"context"
	"log"
	"reflect"
	"time"

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
	Color       *string `json:"color,omitempty"`
}

type CardMoveReq struct {
	ListUID string `json:"list_uid"`
}

type Card struct {
	Title       string     `json:"title"`
	Index       int        `json:"index"`
	UID         string     `json:"uid"`
	Description string     `json:"description"`
	Comments    []Comment  `json:"comments"`
	Files       []FileInfo `json:"files"`
	Timestamp   int64      `json:"timestamp"`
	Color       string     `json:"color"`
}

func (c *CardReq) Save(ctx context.Context, listUID string) (*Card, error) {
	UID, _ := nano.Nanoid()
	newCard := Card{
		UID:         UID,
		Index:       c.Index,
		Title:       c.Title,
		Description: c.Description,
		Comments:    []Comment{},
		Files:       []FileInfo{},
		Timestamp:   time.Now().Unix(),
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

func DeleteCard(c context.Context, cardUID string) error {
	wh := service.DB.Collection(service.ListCollectionName)
	filter := bson.M{"cards": bson.M{"$elemMatch": bson.M{"uid": cardUID}}}
	res, err := wh.UpdateOne(c, filter,
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
	filter := bson.M{"cards": bson.M{"$elemMatch": bson.M{"uid": cardUID}}}
	var toSet bson.D
	combineIfExists(&toSet, "cards.$.index", updateBody.Index)
	combineIfExists(&toSet, "cards.$.title", updateBody.Title)
	combineIfExists(&toSet, "cards.$.description", updateBody.Description)
	combineIfExists(&toSet, "cards.$.color", updateBody.Color)
	update := bson.D{{"$set", toSet}}
	res, err := wh.UpdateOne(c, filter, update)
	if err != nil {
		return err
	} else if res.ModifiedCount == 0 {
		return ResourceNotFound
	}
	return nil
}

func MoveCard(c context.Context, cardUID string, updateBody *CardMoveReq) error {
	log.Printf("List UID: %s\n", updateBody.ListUID)
	card, err := getCard(c, cardUID)
	if err != nil {
		log.Println("Get error")
		log.Println(err)
		return ResourceNotFound
	}
	log.Printf("Card %+v", card)

	err = DeleteCard(c, cardUID)
	if err != nil {
		log.Println("Delete error")
		return err
	}

	lh := service.DB.Collection(service.ListCollectionName)
	_, err = lh.UpdateOne(c, bson.M{"uid": updateBody.ListUID}, bson.D{{"$push", bson.D{{"cards", card}}}})
	if err != nil {
		log.Println("Update error")
		log.Println(err)
		return err
	}

	log.Println("No error")
	return nil
}

func getCard(c context.Context, cardUID string) (*Card, error) {
	wh := service.DB.Collection(service.ListCollectionName)
	var list List
	filter := bson.M{"cards": bson.M{"$elemMatch": bson.M{"uid": cardUID}}}
	err := wh.FindOne(c, filter).Decode(&list)
	if err != nil {
		return nil, err
	}

	for _, v := range list.Cards {
		if v.UID == cardUID {
			return &v, nil
		}
	}
	return nil, ResourceNotFound
}

func combineIfExists(doc *bson.D, key string, val interface{}) {
	if doc == nil {
		doc = &bson.D{}
	}
	if val != nil && !reflect.ValueOf(val).IsNil() {
		*doc = append(*doc, bson.E{Key: key, Value: val})
	}
}
