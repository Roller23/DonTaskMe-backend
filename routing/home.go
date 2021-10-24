package routing

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func homeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, struct{ Status string }{"Operational"})
}
