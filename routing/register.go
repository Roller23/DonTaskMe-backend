package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"github.com/gin-gonic/gin"
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

	err := userReq.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Status(http.StatusCreated)
}
