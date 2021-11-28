package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func getLists(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	boardUID := c.Param("board")
	lists, err := model.FindBoardLists(c, boardUID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, lists)
}

func addList(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	boardUID := c.Param("board")
	var listReq model.ListReq
	err = c.ShouldBindJSON(&listReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// TODO: check that only labradors can add list
	// TODO: add constraint for index (not higher than previous one and all in natural order)
	list, err := listReq.Save(c, boardUID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, list)
}

func deleteList(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	listUID := c.Param("uid")

	//TODO: check if in labradors or something
	err = model.DeleteList(c, listUID)
	if err == model.ResourceNotFound {
		c.JSON(http.StatusBadRequest, err)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	//TODO: send just status
	c.JSON(http.StatusAccepted, "")
}
