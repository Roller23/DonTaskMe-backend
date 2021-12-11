package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"github.com/gin-gonic/gin"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"path/filepath"
)

func saveFile(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "No file is received")
		return
	}

	extension := filepath.Ext(file.Filename)
	UID, _ := nano.Nanoid()
	cardUID := c.Param("card")

	newFilename := UID + extension
	storagePath := "files/" + cardUID + "/"

	fileInfo := model.FileInfo{
		Filename:    newFilename,
		StoragePath: storagePath,
	}

	err = fileInfo.Save(c, cardUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Unable to save the file's metadata")
		return
	}

	if err = c.SaveUploadedFile(file, storagePath+newFilename); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Unable to save the file")
		return
	}

	c.JSON(http.StatusAccepted, "Your file has been successfully uploaded.")
}
