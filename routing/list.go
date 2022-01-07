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

	for i := range lists {
		for j := range lists[i].Cards {
			for k := range lists[i].Cards[j].Files {
				fullPath := fmt.Sprintf("%s/%s", storageUrl, lists[i].Cards[j].Files[k].StoragePath)
				lists[i].Cards[j].Files[k].StoragePath = fullPath
			}
		}
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

	err = model.DeleteList(c, listUID)
	if err == model.ResourceNotFound {
		c.JSON(http.StatusBadRequest, err)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}

func updateList(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var updateReq model.ListUpdateReq
	err = c.ShouldBindJSON(&updateReq)
	if err != nil {
		log.Println("Inappropriate body")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	listUID := c.Param("uid")
	err = model.UpdateList(c, listUID, &updateReq)
	if err == model.ResourceNotFound {
		c.JSON(http.StatusBadRequest, err)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}
