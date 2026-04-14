package handler

import (
	"control_plane/internal/domain"
	"control_plane/internal/service/health"
	dto "control_plane/internal/transport/http_gin/dto/health"
	"log/slog"
	"net/http"
	"time"

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

// @Summary Get client health
// @Tags health
// @Param client_id path string true "Client ID"
// @Success 200 {object} domain.APIClientHealth
// @Router /clients/{client_id}/health [get]
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
		h.log.Warn("health not found → return default",
			"client_id", clientID,
		)

		c.JSON(http.StatusOK, dto.APIClientHealth{
			ClientID:  clientID,
			Status:    domain.HealthUnknown,
			LastCheck: time.Now(),
			Message:   "no health data yet",
		})
		return
	}

	h.log.Info("health fetched",
		"client_id", clientID,
	)

	c.JSON(http.StatusOK, health)
}
