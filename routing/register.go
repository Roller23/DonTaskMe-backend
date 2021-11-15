package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"DonTaskMe-backend/internal/service"
	"DonTaskMe-backend/pkg/hash"
	"context"
	"github.com/gin-gonic/gin"
	nano "github.com/matoous/go-nanoid"
	"net/http"
)

func register(c *gin.Context) {
	var userReq model.UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if helpers.UserAlreadyExists(&userReq.Username) {
		c.JSON(http.StatusConflict, "userReq with given username already exists")
		return
	}

	//TODO: validate password
	//if !isPasswordValid(ld.Password) {
	//	c.JSON(
	//		http.StatusNotAcceptable,
	//		errorMsg{ Message: "Password does not meet minimal rules: eight characters, one digit, one capital letter, one special character, one lowercase letter" },
	//	)
	//}

	hashedPassword, err := hash.Generate(&userReq.Password)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, err.Error())
		return
	}
	userReq.Password = hashedPassword

	uid, err := nano.Nanoid()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	newUser := model.User{
		Uid:      &uid,
		Username: userReq.Username,
		Password: userReq.Password,
		Token:    nil,
	}

	usersCollection := service.Client.Database(service.Name).Collection(service.UsersCollectionName)
	_, err = usersCollection.InsertOne(context.TODO(), newUser)
	c.Status(http.StatusCreated)
}
