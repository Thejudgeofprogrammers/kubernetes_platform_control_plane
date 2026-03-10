package handler

import (
	"control_plane/internal/service/health"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	service health.HealthService
}

func NewHealthHandler(s health.HealthService) *HealthHandler {
	return &HealthHandler{
		service: s,
	}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	clientID := c.Param("client_id")

	health, err := h.service.Get(
		c.Request.Context(),
		clientID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch health",
		})
		return
	}

	if health == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "health not found",
		})
		return
	}
	
	c.JSON(http.StatusOK, health)
}
