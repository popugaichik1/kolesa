package auth_transport_http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	core_errors "github.com/zosinkin/social_network/internal/core/errors"
)

type CreateUserResponse UserDTOResponse

func (h *AuthHTTPHandler) Register(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"failed to decode and validate HTTP request": err.Error(),
		})
		return
	}

	user, err := h.authService.Register(
		c.Request.Context(),
		req.Username,
		req.PhoneNumber,
		req.Password,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"failed to create user": err.Error(),
		})
	}
	resp := UserDTOResponse{
		ID:          user.ID,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
	}

	c.JSON(http.StatusCreated, resp)
}


func (h *AuthHTTPHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid request payload",
		})
		c.Abort()
		return
	}

	refreshTokenTTL := 7 * 24 * time.Hour

	accessToken, refreshToken, err := h.authService.LoginWithRefresh(
		c.Request.Context(),
		req.PhoneNumber,
		req.Password,
		refreshTokenTTL,
	)
	if err != nil {
		if errors.Is(err, core_errors.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}
	}

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusOK, response)
}


