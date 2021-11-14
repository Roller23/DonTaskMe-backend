package routing

import (
	"DonTaskMe-backend/internal/db"
	"DonTaskMe-backend/internal/models"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

type WorkspacesRequest struct {
	Token string `json:"token"`
}

func getWorkspaces(c *gin.Context) {
	wh := db.Handler.Collection(db.WorkspaceCollectionName)
	var token WorkspacesRequest
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := findUserByID(token.Token)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, "user does not exists in database")
		return
	}

	var workspaces []models.Workspace
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
	//wh := db.Handler.Collection(db.WorkspaceCollectionName)

}

func findUserByID(token string) (*models.User, error) {
	var res models.User
	usersCollection := db.Handler.Collection(db.UsersCollectionName)
	err := usersCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
