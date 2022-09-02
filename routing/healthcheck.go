package routing

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, struct{ Status string }{"Operational"})
}
