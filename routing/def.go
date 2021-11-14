package routing

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetServer(mode string) *gin.Engine {
	gin.SetMode(mode)
	router := gin.Default()
	router.Use(cors.Default())

	//Set up routes here
	router.GET("/health-check", healthCheckHandler)
	router.POST("/register", register)
	router.POST("/login", login)
	router.POST("/workspaces", getWorkspaces) //TODO change to GET someday
	return router
}
