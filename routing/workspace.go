package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"DonTaskMe-backend/internal/service"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	workspace, err := workspaceReq.Save(*user.Uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, workspace)
}

func updateWorkspace(c *gin.Context) {
	panic("Not implemented yet!")
}

func deleteWorkspace(c *gin.Context) {
	UID := c.Param("uid")
	err := model.Delete(UID)
	if err == model.ResourceNotFound {
		c.JSON(http.StatusBadRequest, err)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	//TODO: send just status
	c.JSON(http.StatusAccepted, "")
}
