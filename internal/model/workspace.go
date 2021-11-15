package model

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/service"
	"context"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
)

type WorkspaceRequest struct {
	Token     string   `json:"token"`
	Title     string   `json:"title"`
	Desc      string   `json:"desc"`
	Boards    []Board  `json:"boards,omitempty"`
	Labradors []string `json:"labradors,omitempty"`
}

type Workspace struct {
	UID       string   `json:"uid"`
	Title     string   `json:"title"`
	Desc      string   `json:"desc"`
	Boards    []Board  `json:"boards"`
	Owner     string   `json:"owner"`
	Labradors []string `json:"labradors"` //users ID
}

func (w *WorkspaceRequest) Save(ownerUID string) error {
	UID, _ := nano.Nanoid()
	newWorkspace := Workspace{
		UID:       UID,
		Title:     w.Title,
		Desc:      w.Desc,
		Boards:    w.Boards,
		Owner:     ownerUID,
		Labradors: w.Labradors,
	}

	wh := service.DB.Collection(service.WorkspaceCollectionName)
	_, err := wh.InsertOne(context.TODO(), newWorkspace)
	return err
}

func Delete(workspaceUID string) error {
	wh := service.DB.Collection(service.WorkspaceCollectionName)
	res, err := wh.DeleteOne(context.TODO(), bson.D{{"uid", workspaceUID}})
	if err != nil {
		return err
	} else if res.DeletedCount == 0 {
		return helpers.ResourceNotFound
	}
	return nil
}
