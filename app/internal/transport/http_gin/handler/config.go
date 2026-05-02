package handler

import (
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	cfgService "control_plane/internal/service/config"
	configDTO "control_plane/internal/transport/http_gin/dto/config"
	"control_plane/internal/transport/http_gin/mapper"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	service cfgService.ConfigService
	log     logger.Logger
}

func NewConfigHandler(s cfgService.ConfigService, log logger.Logger) *ConfigHandler {
	return &ConfigHandler{
		service: s,
		log:     log,
	}
}

// @Summary List configs
// @Tags configs
// @Param client_id path string true "Client ID"
// @Success 200 {array} dto.ConfigResponse
// @Failure 404 {object} map[string]string
// @Router /clients/{client_id}/configs [get]
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

	resp := make([]configDTO.ConfigResponse, 0, len(configs))
	for _, cfg := range configs {
		resp = append(resp, mapper.ToConfigResponse(cfg))
	}

	h.log.Info("configs listed",
		"client_id", clientID,
		"count", len(resp),
	)

	c.JSON(http.StatusOK, resp)
}

// @Summary Create config
// @Tags configs
// @Param client_id path string true "Client ID"
// @Accept json
// @Param request body dto.ClientConfigRequest true "Config data"
// @Success 201 {object} dto.ConfigResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /clients/{client_id}/configs [post]
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

	var req configDTO.ClientConfigRequest
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

// @Summary Deploy config
// @Description Активирует конфигурацию и делает rolling update
// @Tags configs
// @Param client_id path string true "Client ID"
// @Param config_id path string true "Config ID"
// @Success 202 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /clients/{client_id}/configs/{config_id}/deploy [post]
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

// @Summary Delete config
// @Tags configs
// @Param client_id path string true "Client ID"
// @Param config_id path string true "Config ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /clients/{client_id}/configs/{config_id}/delete [delete]
func (h *ConfigHandler) Delete(c *gin.Context) {
	clientID := c.Param("client_id")
	configID := c.Param("config_id")

	h.log.Info("http delete config started",
		"client_id", clientID,
		"config_id", configID,
	)

	if clientID == "" || configID == "" {
		h.log.Warn("missing client_id or config_id")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id and config id are required",
		})
		return
	}

	err := h.service.Delete(
		c.Request.Context(),
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
		case errors.Is(err, domain.ErrConfigNotFound):
			h.log.Warn("config not found",
				"config_id", configID,
			)

			c.JSON(http.StatusNotFound, gin.H{"error": "config not found"})
		case errors.Is(err, domain.ErrInvalidStateTransition):
			h.log.Warn("cannot delete active config",
				"client_id", clientID,
				"config_id", configID,
			)

			c.JSON(http.StatusConflict, gin.H{"error": "cannot delete active config"})
		default:
			h.log.Error("failed to delete config",
				"client_id", clientID,
				"config_id", configID,
				"error", err,
			)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete config"})
		}
		return
	}

	h.log.Info("config deleted",
		"client_id", clientID,
		"config_id", configID,
	)

	c.Status(http.StatusNoContent)
}

// @Summary Update config
// @Tags configs
// @Param client_id path string true "Client ID"
// @Param config_id path string true "Config ID"
// @Accept json
// @Param request body dto.ClientConfigRequest true "Config data"
// @Success 200 {object} dto.ConfigResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clients/{client_id}/configs/{config_id}/update [put]
func (h *ConfigHandler) Update(c *gin.Context) {
	clientID := c.Param("client_id")
	configID := c.Param("config_id")
	userID := c.GetString("user_id")

	h.log.Info("http update config started",
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

	var req configDTO.ClientConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid update config body",
			"client_id", clientID,
			"config_id", configID,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	config, err := h.service.Update(
		c.Request.Context(),
		userID,
		clientID,
		configID,
		req,
	)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrClientNotFound):
			h.log.Warn("client not found",
				"client_id", clientID,
			)

			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		case errors.Is(err, domain.ErrConfigNotFound):
			h.log.Warn("config not found",
				"config_id", configID,
			)

			c.JSON(http.StatusNotFound, gin.H{"error": "config not found"})
		default:
			h.log.Error("failed to update config",
				"client_id", clientID,
				"config_id", configID,
				"user_id", userID,
				"error", err,
			)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update config"})
		}
		return
	}

	h.log.Info("config updated",
		"client_id", clientID,
		"config_id", configID,
		"user_id", userID,
	)

	c.JSON(http.StatusOK, mapper.ToConfigResponse(config))
}
