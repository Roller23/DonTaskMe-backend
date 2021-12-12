package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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

	listUID := c.Query("list")
	cards, err := model.FindListCards(c, listUID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for _, card := range cards {
		for _, file := range card.Files {
			file.StoragePath = fmt.Sprintf("%s/%s", storageUrl, file.StoragePath)
		}
	}

	c.JSON(http.StatusOK, cards)
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
		c.Writer.WriteHeader(http.StatusBadRequest)
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

	card, err := cardReq.Save(c, list.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, card)
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

	var updateReq model.CardUpdateReq
	err = c.ShouldBindJSON(&updateReq)
	if err != nil {
		log.Println("Inappropriate body")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	cardUID := c.Param("card")
	err = model.UpdateCard(c, cardUID, updateReq)

	if err == model.ResourceNotFound {
		log.Println("No such resource")
		c.JSON(http.StatusBadRequest, err)
		return
	} else if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Writer.WriteHeader(http.StatusAccepted)
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

	cardUID := c.Param("card")
	err = model.DeleteCard(c, cardUID)
	if err == model.ResourceNotFound {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}
