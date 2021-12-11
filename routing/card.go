package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
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

	//TODO: check if workspace labrador in

	listUID := c.Param("list")
	// _, err = model.FindList(c, listUID)
	// if err == mongo.ErrNoDocuments {
	// 	c.JSON(http.StatusBadRequest, "no such list")
	// 	return
	// } else if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err.Error())
	// 	return
	// }

	cardUID := c.Param("card")
	err = model.DeleteCard(c, listUID, cardUID)
	if err == model.ResourceNotFound {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}
