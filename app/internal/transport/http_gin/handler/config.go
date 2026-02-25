package handler

import (
	"control_plane/internal/domain"
	"control_plane/internal/service"
	dto "control_plane/internal/transport/http_gin/dto/client"
	"control_plane/internal/transport/http_gin/mapper"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	service service.ConfigService
}

func NewConfigHandler(s service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		service: s,
	}
}

func (h *ConfigHandler) ListConfigs(c *gin.Context) {
	clientID := c.Param("client_id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id is required",
		})
		return
	}

	configs, err := h.service.ListConfigs(c.Request.Context(), clientID)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "client not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list configs",
		})
		return
	}

	resp := make([]dto.ConfigResponse, 0, len(configs))
	for _, cfg := range configs {
		resp = append(resp, mapper.ToConfigResponse(cfg))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ConfigHandler) CreateConfig(c *gin.Context) {
	clientID := c.Param("client_id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id is required",
		})
		return
	}

	var req dto.ClientConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userID := c.GetString("user_id")

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
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		case errors.Is(err, domain.ErrConfigVersionExists):
			c.JSON(http.StatusConflict, gin.H{"error": "config version already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create config"})
		}
		return
	}

	c.JSON(http.StatusCreated, mapper.ToConfigResponse(config))
}

func (h *ConfigHandler) Deploy(c *gin.Context) {
	clientID := c.Param("client_id")
	configID := c.Param("config_id")

	if clientID == "" || configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id and config id are required",
		})
		return
	}

	userID := c.GetString("user_id")

	err := h.service.Deploy(
		c.Request.Context(),
		userID,
		clientID,
		configID,
	)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrClientNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		case errors.Is(err, domain.ErrInvalidStateTransition):
			c.JSON(http.StatusConflict, gin.H{"error": "invalid state transaction"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deploy config"})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"client_id": clientID,
		"config_id": configID,
		"status": string(domain.ClientStatusDeploying),
	})
}
