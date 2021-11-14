package routing

import (
	db2 "DonTaskMe-backend/internal/db"
	models2 "DonTaskMe-backend/internal/models"
	"DonTaskMe-backend/pkg/hash"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
)

func login(c *gin.Context) {
	var user models2.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dbUser, err := findUser(&user.Username)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, "user does not exists in database")
		return
	}

	if !hash.CheckPasswordHash(&user.Password, &dbUser.Password) {
		c.JSON(http.StatusNotAcceptable, "invalid login or password")
		return
	}

	c.JSON(http.StatusOK, "correct")
}

func findUser(username *string) (*models2.User, error) {
	var res models2.User
	usersCollection := db2.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("USERS_COLLECTION"))
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
