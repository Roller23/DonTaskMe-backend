package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"DonTaskMe-backend/pkg/hash"
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Token string `json:"token"`
}

func login(c *gin.Context) {
	var user model.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dbUser, err := helpers.FindUser(&user.Username)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, "user does not exists in database")
		return
	}

	if !hash.CheckPasswordHash(&user.Password, &dbUser.Password) {
		c.JSON(http.StatusNotAcceptable, "invalid login or password")
		return
	}

	err = dbUser.AssignNewToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "could not create token")
		return
	}

	c.JSON(http.StatusOK, response{Token: *dbUser.Token})
}
