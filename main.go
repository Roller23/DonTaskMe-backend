package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Cheers!")
	})
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct{ Status string }{"Operational"})
	})
	err := router.Run()
	if err != nil {
		log.Fatalln("Could not start the server:", err)
	}
}
