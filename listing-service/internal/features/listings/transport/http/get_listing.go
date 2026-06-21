package listings_transport_http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	core_errors "listing-service/internal/core/errors"
)

func (h *ListingsHandler) GetListing(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid listing id"})
		return
	}

	listing, err := h.service.GetListingByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, toListingResponse(listing))
}
