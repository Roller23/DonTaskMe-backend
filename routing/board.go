package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func getBoards(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	workspaceUID := c.Param("workspace")
	boards, err := model.FindWorkspaceBoards(c, workspaceUID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, boards)
}

func addBoard(c *gin.Context) {
	token := c.Query("token")
	user, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var boardReq model.BoardRequest
	err = c.ShouldBindJSON(&boardReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	workspace, err := model.FindWorkspace(c, boardReq.WorkspaceUID)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, "no such workspace")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if *user.Uid != workspace.Owner {
		c.JSON(http.StatusExpectationFailed, "no owner privilege")
		return
	}

	board, err := boardReq.Save(c, workspace.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, board)
}

func updateBoard(c *gin.Context) {
	panic("Not implemented yet!")
}

func deleteBoard(c *gin.Context) {
	token := c.Query("token")
	user, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	workspaceUID := c.Param("workspace")
	workspace, err := model.FindWorkspace(c, workspaceUID)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, "no such workspace")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	boardUID := c.Param("board")
	if workspace.Owner == *user.Uid {
		err = model.DeleteBoard(c, workspaceUID, boardUID)
		if err == model.ResourceNotFound {
			c.JSON(http.StatusBadRequest, err)
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		//TODO: send just status
		c.JSON(http.StatusAccepted, "")
		return
	}

	c.JSON(http.StatusBadRequest, "no ownership")
}
