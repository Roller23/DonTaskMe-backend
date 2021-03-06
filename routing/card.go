package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"fmt"
	"log"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func isHexColor(str string) bool {
	if len(str) > 7 || len(str) < 4 {
		return false
	}
	if str[0] != '#' {
		return false
	}
	numbers := str[1:]
	for _, char := range numbers {
		if !(unicode.IsDigit(char) || unicode.IsLetter(char)) {
			return false
		}
	}
	return true
}

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

	for i := range cards {
		for j := range cards[i].Files {
			cards[i].Files[j].StoragePath = fmt.Sprintf("%s/%s", storageUrl, cards[i].Files[j].StoragePath)
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

	if updateReq.Color != nil && !isHexColor(*updateReq.Color) {
		log.Println(*updateReq.Color, "is not a hex color")
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

func moveCard(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var updateReq model.CardMoveReq
	err = c.ShouldBindJSON(&updateReq)
	if err != nil {
		log.Println("Inappropriate body")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	cardUID := c.Param("card")
	err = model.MoveCard(c, cardUID, &updateReq)

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
