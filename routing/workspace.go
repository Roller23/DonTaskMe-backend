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
	if token == "" {
		c.JSON(http.StatusExpectationFailed, "no token passed")
		return
	}
	wh := service.DB.Collection(service.WorkspaceCollectionName)

	user, err := helpers.FindUserByToken(token)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, "user does not exists in database")
		return
	}

	workspaces := make([]model.Workspace, 0)
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
	//wh := db.DB.Collection(db.WorkspaceCollectionName)

}
