package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"DonTaskMe-backend/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

func getWorkspaces(c *gin.Context) {
	token := c.Query("token")
	user, err := helpers.FindUserByToken(token)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, "access denied")
		return
	}

	workspaces := make([]model.Workspace, 0)
	wh := service.DB.Collection(service.WorkspaceCollectionName)
	cursor, err := wh.Find(context.TODO(), bson.M{"labradors": bson.M{"$in": []string{*user.Uid}}})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = cursor.All(context.TODO(), &workspaces)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

func addWorkspace(c *gin.Context) {
	var workspaceReq model.WorkspaceRequest
	err := c.ShouldBindJSON(&workspaceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := helpers.FindUserByToken(workspaceReq.Token)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, "access denied")
		return
	}

	wh := service.DB.Collection(service.WorkspaceCollectionName)
	newWorkspace := model.Workspace{
		Title:     workspaceReq.Title,
		Desc:      workspaceReq.Desc,
		Boards:    workspaceReq.Boards,
		Labradors: append(workspaceReq.Labradors, *user.Uid),
	}

	_, err = wh.InsertOne(context.TODO(), newWorkspace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
}
