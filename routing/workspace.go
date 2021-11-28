package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"DonTaskMe-backend/internal/service"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func getWorkspaces(c *gin.Context) {
	token := c.Query("token")
	user, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	workspaces, err := model.FindUsersWorkspaces(c, *user.Uid, false)
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

	workspace, err := workspaceReq.Save(c, *user.Uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, workspace)
}

func updateWorkspace(_ *gin.Context) {
	panic("Not implemented yet!")
}

func deleteWorkspace(c *gin.Context) {
	token := c.Query("token")
	user, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	workspaces, err := model.FindUsersWorkspaces(c, *user.Uid, true)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	UID := c.Param("uid")
	for _, w := range workspaces {
		if w.UID == UID {
			err = model.Delete(c, UID)
			if err == model.ResourceNotFound {
				c.JSON(http.StatusBadRequest, err)
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
			//TODO: send just status
			c.JSON(http.StatusAccepted, "")
			return
		}
	}

	c.JSON(http.StatusBadRequest, "no ownership")
}
