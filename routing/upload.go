package routing

import (
	"DonTaskMe-backend/internal/helpers"
	"DonTaskMe-backend/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	storageUrl = "http://dontstoreme.vxm.pl"
)

type FileResponse struct {
	Success bool   `json:"success,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Path    string `json:"path,omitempty"`
}

func UploadFile(path string) (FileResponse, error) {
	file, err := os.Open(path)
	if err != nil {
		return FileResponse{}, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	part2, _ := writer.CreateFormField("key")
	io.Copy(part2, bytes.NewBufferString("97e5fe51-dcd4-4bdf-ac2e-0d96c990f9fc"))
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, storageUrl, body)
	if err != nil {
		return FileResponse{}, err
	}

	client := &http.Client{}
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Content-Type", writer.FormDataContentType())

	log.Println("uploading", filepath.Base(file.Name()))

	resp, err := client.Do(req)
	if err != nil {
		return FileResponse{}, err
	}

	log.Println("upload completed")

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return FileResponse{}, err
	}

	var response FileResponse
	err = json.Unmarshal(content, &response)
	return response, err
}

func saveFile(c *gin.Context) {
	token := c.PostForm("token")
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

	newFilename := UID + extension
	storagePath := "uploads/"
	newPath := storagePath + newFilename

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		os.Mkdir(storagePath, 0777)
	}

	if err = c.SaveUploadedFile(file, newPath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Unable to save the file")
		return
	}

	res, err := UploadFile(newPath)

	if err != nil {
		log.Println("err", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "File upload failed")
		return
	}

	if !res.Success {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "File upload failed: "+res.Msg)
		return
	}

	cardUID := c.Param("card")
	fileUID, _ := nano.Nanoid()
	fileInfo := model.FileInfo{
		Filename:    file.Filename,
		StoragePath: res.Path,
		UID:         fileUID,
	}

	err = fileInfo.Save(c, cardUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Unable to save the file's metadata")
		return
	}

	fileInfo.StoragePath = fmt.Sprintf("%s/%s", storageUrl, fileInfo.StoragePath)
	c.JSON(http.StatusAccepted, fileInfo)
}

func deleteFile(c *gin.Context) {
	token := c.Query("token")
	_, err := helpers.FindUserByToken(token)
	if err == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusExpectationFailed, err)
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	cardUID := c.Param("card")
	fileUID := c.Param("file")

	err = model.DeleteFile(c, cardUID, fileUID)
	if err == model.ResourceNotFound {
		c.AbortWithStatusJSON(http.StatusBadRequest, "No such file on the server")
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}
