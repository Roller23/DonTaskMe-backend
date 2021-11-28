package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func getCards(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	listUID := c.Param("list")
	boards, err := model.FindListCards(c, listUID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, boards)
}

func addCard(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var cardReq model.CardReq
	err = c.ShouldBindJSON(&cardReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	list, err := model.FindList(c, cardReq.ListUID)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, "no such list")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	//TODO: check if list labradors contains user uid

	board, err := cardReq.Save(c, list.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, board)
}

func updateCard(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	//TODO: update
	deleteCard(c)
	addCard(c)
}

func deleteCard(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	listUID := c.Param("list")
	_, err = model.FindList(c, listUID)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, "no such list")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	//TODO: check if workspace labrador in

	cardUID := c.Param("cardUID")
	err = model.DeleteCard(c, listUID, cardUID)
	if err == model.ResourceNotFound {
		c.JSON(http.StatusBadRequest, err)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	//TODO: send just status
	c.JSON(http.StatusAccepted, "")
}
