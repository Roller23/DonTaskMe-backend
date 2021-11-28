package model

import (
	"DonTaskMe-backend/internal/service"
	"context"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
)

type BoardRequest struct {
	WorkspaceUID string `json:"workspace"`
	Title        string `json:"title"`
}

type Board struct {
	Title string `json:"title"`
	UID   string `json:"uid"`
}

func (w *BoardRequest) Save(c context.Context, workspaceUID string) (*Board, error) {
	UID, _ := nano.Nanoid()
	newBoard := Board{
		UID:   UID,
		Title: w.Title,
	}

	wh := service.DB.Collection(service.WorkspaceCollectionName)
	_, err := wh.UpdateOne(c, bson.M{"uid": workspaceUID}, bson.D{{"$push", bson.D{{"boards", newBoard}}}})
	if err != nil {
		return nil, err
	}
	return &newBoard, err
}

func FindWorkspaceBoards(c context.Context, workspaceUID string) ([]Board, error) {
	wh := service.DB.Collection(service.WorkspaceCollectionName)
	var workspace Workspace
	err := wh.FindOne(c, bson.D{{"uid", workspaceUID}}).Decode(&workspace)
	if err != nil {
		return nil, err
	}

	return workspace.Boards, nil
}

func DeleteBoard(c context.Context, workspaceUID string, boardUID string) error {
	wh := service.DB.Collection(service.WorkspaceCollectionName)
	res, err := wh.UpdateOne(
		c, bson.D{{"uid", workspaceUID}},
		bson.D{{"$pull", bson.D{{"boards", bson.D{{"uid", boardUID}}}}}},
	)

	if err != nil {
		return err
	} else if res.ModifiedCount == 0 {
		return ResourceNotFound
	}
	return nil
}
