package auth_transport_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



func (h *AuthHTTPHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Health": "ok",
	})
	return
}


func (h *AuthHTTPHandler) Authorized(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Authorized": "OK",
	})
	return
}