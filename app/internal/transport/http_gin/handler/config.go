package handler

import (
	"control_plane/internal/domain"
	cfgService "control_plane/internal/service/config"
	dto "control_plane/internal/transport/http_gin/dto/client"
	"control_plane/internal/transport/http_gin/mapper"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	service cfgService.ConfigService
	log     *slog.Logger
}

func NewConfigHandler(s cfgService.ConfigService, log *slog.Logger) *ConfigHandler {
	return &ConfigHandler{
		service: s,
		log:     log,
	}
}

func (h *ConfigHandler) ListConfigs(c *gin.Context) {
	clientID := c.Param("client_id")

	h.log.Info("http list configs started",
		"client_id", clientID,
	)

	if clientID == "" {
		h.log.Warn("missing client_id")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id is required",
		})
		return
	}

	configs, err := h.service.ListConfigs(c.Request.Context(), clientID)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			h.log.Warn("client not found",
				"client_id", clientID,
			)

			c.JSON(http.StatusNotFound, gin.H{
				"error": "client not found",
			})
			return
		}
		h.log.Error("failed to list configs",
			"client_id", clientID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list configs",
		})
		return
	}

	resp := make([]dto.ConfigResponse, 0, len(configs))
	for _, cfg := range configs {
		resp = append(resp, mapper.ToConfigResponse(cfg))
	}

	h.log.Info("configs listed",
		"client_id", clientID,
		"count", len(resp),
	)

	c.JSON(http.StatusOK, resp)
}

func (h *ConfigHandler) CreateConfig(c *gin.Context) {
	clientID := c.Param("client_id")
	userID := c.GetString("user_id")

	h.log.Info("http create config started",
		"client_id", clientID,
		"user_id", userID,
	)

	if clientID == "" {
		h.log.Warn("missing client_id")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id is required",
		})
		return
	}

	var req dto.ClientConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid create config body",
			"client_id", clientID,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	config, err := h.service.CreateConfig(
		c.Request.Context(),
		userID,
		clientID,
		req.Version,
		domain.AuthType(req.AuthType),
		req.AuthRef,
		req.TimeoutMs,
		req.RetryCount,
		req.RetryBackoff,
		req.Headers,
	)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrClientNotFound):
			h.log.Warn("client not found",
				"client_id", clientID,
			)

			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		case errors.Is(err, domain.ErrConfigVersionExists):
			h.log.Warn("config version exists",
				"client_id", clientID,
				"version", req.Version,
			)

			c.JSON(http.StatusConflict, gin.H{"error": "config version already exists"})
		default:
			h.log.Error("failed to create config",
				"client_id", clientID,
				"user_id", userID,
				"error", err,
			)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create config"})
		}
		return
	}

	h.log.Info("config created",
		"client_id", clientID,
		"config_id", config.ID,
		"user_id", userID,
		"version", config.Version,
	)

	c.JSON(http.StatusCreated, mapper.ToConfigResponse(config))
}

func (h *ConfigHandler) Deploy(c *gin.Context) {
	clientID := c.Param("client_id")
	configID := c.Param("config_id")
	userID := c.GetString("user_id")

	h.log.Info("http deploy config started",
		"client_id", clientID,
		"config_id", configID,
		"user_id", userID,
	)

	if clientID == "" || configID == "" {
		h.log.Warn("missing client_id or config_id")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id and config id are required",
		})
		return
	}

	err := h.service.Deploy(
		c.Request.Context(),
		userID,
		clientID,
		configID,
	)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrClientNotFound):
			h.log.Warn("client not found",
				"client_id", clientID,
			)

			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		case errors.Is(err, domain.ErrInvalidStateTransition):
			h.log.Warn("invalid state transition",
				"client_id", clientID,
				"config_id", configID,
			)

			c.JSON(http.StatusConflict, gin.H{"error": "invalid state transaction"})
		default:
			h.log.Error("deploy config failed",
				"client_id", clientID,
				"config_id", configID,
				"user_id", userID,
				"error", err,
			)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deploy config"})
		}
		return
	}

	h.log.Info("deploy scheduled",
		"client_id", clientID,
		"config_id", configID,
		"user_id", userID,
	)

	c.JSON(http.StatusAccepted, gin.H{
		"client_id": clientID,
		"config_id": configID,
		"status":    string(domain.ClientStatusDeploying),
	})
}
