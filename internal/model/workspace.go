package model

import (
	"DonTaskMe-backend/internal/service"
	"context"
	"errors"

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

var (
	ResourceNotFound = errors.New("no such resource")
)

func (w *WorkspaceRequest) Save(c context.Context, ownerUID string) (*Workspace, error) {
	UID, _ := nano.Nanoid()
	labradors := []string{ownerUID}
	newWorkspace := Workspace{
		UID:       UID,
		Title:     w.Title,
		Desc:      w.Desc,
		Boards:    []Board{},
		Owner:     ownerUID,
		Labradors: labradors,
	}

	wh := service.DB.Collection(service.WorkspaceCollectionName)
	_, err := wh.InsertOne(c, newWorkspace)
	return &newWorkspace, err
}

func Delete(c context.Context, workspaceUID string) error {
	wh := service.DB.Collection(service.WorkspaceCollectionName)
	res, err := wh.DeleteOne(c, bson.D{{"uid", workspaceUID}})
	if err != nil {
		return err
	} else if res.DeletedCount == 0 {
		return ResourceNotFound
	}
	return nil
}

func FindUsersWorkspaces(c context.Context, userUID string, owner bool) ([]Workspace, error) {
	var accessLevel string
	if owner {
		accessLevel = "owner"
	} else {
		accessLevel = "labradors"
	}
	workspaces := make([]Workspace, 0)
	wh := service.DB.Collection(service.WorkspaceCollectionName)
	cursor, err := wh.Find(c, bson.M{accessLevel: bson.M{"$in": []string{userUID}}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(c, &workspaces)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}
