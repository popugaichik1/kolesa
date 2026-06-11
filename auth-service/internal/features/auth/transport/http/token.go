package auth_transport_http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	core_errors "github.com/zosinkin/social_network/internal/core/errors"
)



func (h *AuthHTTPHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid request payload": err.Error(),
		})
		return
	}

	token, err := h.authService.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, core_errors.ErrInvalidToken) || errors.Is(err, core_errors.ErrExpiredToken) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Invalid or expired refresh token": err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Internal server error": err.Error(),
			})
		}
		return
	}

	response := RefreshResponse{Token: token}

	c.JSON(http.StatusOK, response)
}
