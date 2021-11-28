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

	router.GET("/workspaces", getWorkspaces)
	router.POST("/workspaces", addWorkspace)
	router.PATCH("/workspaces/:uid", updateWorkspace)
	router.DELETE("/workspaces/:uid", deleteWorkspace)

	router.GET("/boards/:workspace", getBoards)
	router.POST("/boards", addBoard)
	router.PATCH("/boards/:uid", updateBoard)
	router.DELETE("/boards/:workspace/:board", deleteBoard)

	return router
}
