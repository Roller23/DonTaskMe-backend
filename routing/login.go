package routing

import (
	"DonTaskMe-backend/internal/db"
	"DonTaskMe-backend/internal/models"
	"DonTaskMe-backend/pkg/hash"
	"context"
	"github.com/gin-gonic/gin"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type response struct {
	Token string `json:"token"`
}

func login(c *gin.Context) {
	var user models.User
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

	err = assignNewToken(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "could not create token")
		return
	}

	c.JSON(http.StatusOK, response{Token: *dbUser.Token})
}

func findUser(username *string) (*models.User, error) {
	var res models.User
	usersCollection := db.Handler.Collection(db.UsersCollectionName)
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func assignNewToken(user *models.User) error {
	uid, _ := nano.Nanoid()
	user.Token = &uid
	usersCollection := db.Handler.Collection(db.UsersCollectionName)
	_, err := usersCollection.UpdateOne(context.TODO(), bson.M{"uid": user.Uid}, user)
	return err
}
