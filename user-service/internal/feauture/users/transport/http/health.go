package transport_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



func (h *HTTPHandler) HealthCheck (c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Health": "ok",
	})
	return
}