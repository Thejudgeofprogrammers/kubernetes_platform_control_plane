package handler

import (
	"control_plane/internal/service/health"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	service health.HealthService
	log     *slog.Logger
}

func NewHealthHandler(s health.HealthService, log *slog.Logger) *HealthHandler {
	return &HealthHandler{
		service: s,
		log:     log,
	}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	clientID := c.Param("client_id")

	h.log.Info("http get health started",
		"client_id", clientID,
	)

	health, err := h.service.Get(
		c.Request.Context(),
		clientID,
	)

	if err != nil {
		h.log.Error("get health failed",
			"client_id", clientID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch health",
		})
		return
	}

	if health == nil {
		h.log.Warn("health not found",
			"client_id", clientID,
		)

		c.JSON(http.StatusNotFound, gin.H{
			"error": "health not found",
		})
		return
	}

	h.log.Info("health fetched",
		"client_id", clientID,
	)

	c.JSON(http.StatusOK, health)
}
