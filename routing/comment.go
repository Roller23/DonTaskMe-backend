package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func addComment(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var commentReq model.CommentReq
	err = c.ShouldBindJSON(&commentReq)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	comment, err := commentReq.Save(c, c.Param("card"))
	if err == model.ResourceNotFound {
		c.JSON(http.StatusBadRequest, "Resource not found")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, comment)
}
