package transport_http

import (
	"context"
	core_logger "user-service/internal/core/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HTTPHandler struct {
	service Service
	log  	*core_logger.Logger
}

func NewHTTPHandler(
	userRepo Service,
) *HTTPHandler {
	return &HTTPHandler{
		service: userRepo,
	}
}

type Service interface {
	SaveUser(
		ctx context.Context,
		ID uuid.UUID,
		username string,
		phoneNumber string, 
	) (error)
}


func (h *HTTPHandler) InitRoutes() *gin.Engine {
	router := gin.Default()

	routes := router.Group("/api/user")
	{
		routes.GET("/health", h.HealthCheck)
		//routes.POST("/save", h.SaveUser)
	}
	return router
} 