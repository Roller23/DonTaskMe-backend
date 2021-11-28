package model

import (
	"DonTaskMe-backend/internal/service"
	"context"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
)

type ListReq struct {
	Title string `json:"title"`
	Index int    `json:"index"`
}

type List struct {
	UID      string `json:"uid"`
	Title    string `json:"title"`
	Index    int    `json:"index"`
	Cards    []Card `json:"cards"`
	BoardUID string `json:"boardUid"`
}

func (l *ListReq) Save(c context.Context, boardUID string) (*List, error) {
	UID, _ := nano.Nanoid()
	newList := List{
		UID:      UID,
		Title:    l.Title,
		Cards:    []Card{},
		Index:    l.Index,
		BoardUID: boardUID,
	}

	lh := service.DB.Collection(service.ListCollectionName)
	_, err := lh.InsertOne(c, newList)

	if err != nil {
		return nil, err
	}

	return &newList, nil
}

func FindBoardLists(c context.Context, boardUID string) ([]List, error) {
	lists := make([]List, 0)
	lh := service.DB.Collection(service.ListCollectionName)
	cursor, err := lh.Find(c, bson.M{"boarduid": boardUID})
	if err != nil {
		return nil, err
	}

	err = cursor.All(c, &lists)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func DeleteList(c context.Context, listUID string) error {
	lh := service.DB.Collection(service.ListCollectionName)
	res, err := lh.DeleteOne(c, bson.D{{"uid", listUID}})
	if err != nil {
		return err
	} else if res.DeletedCount == 0 {
		return ResourceNotFound
	}
	return nil
}

func FindList(c context.Context, listUID string) (*List, error) {
	wh := service.DB.Collection(service.ListCollectionName)
	var list List
	err := wh.FindOne(c, bson.D{{"uid", listUID}}).Decode(&list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}
