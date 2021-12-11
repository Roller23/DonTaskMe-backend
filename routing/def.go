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
	router.PUT("/workspaces/:uid", updateWorkspace)
	router.DELETE("/workspaces/:uid", deleteWorkspace)

	router.GET("/boards/:workspace", getBoards)
	router.POST("/boards", addBoard)
	router.PUT("/boards/:uid", updateBoard)
	router.DELETE("/boards/:workspace/:board", deleteBoard)

	router.GET("/lists/:board", getLists)
	router.POST("/lists/:board", addList)
	router.DELETE("/lists/:uid", deleteList)

	router.GET("/cards", getCards)
	router.POST("/cards", addCard)
	router.PUT("/cards/:card", updateCard)
	router.DELETE("/cards/:list/:card", deleteCard)

	router.POST("cards/:card/comment", addComment)

	return router
}
