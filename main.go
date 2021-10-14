package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Cheers!")
	})
	err := router.Run()
	if err != nil {
		println(fmt.Errorf("%s", "Could not start the server."))
	}
}
