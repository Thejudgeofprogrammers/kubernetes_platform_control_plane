package handler

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/service/client"
	dto "control_plane/internal/transport/http_gin/dto/client"
	"control_plane/internal/transport/http_gin/mapper"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	service client.ClientService
	log     logger.Logger
	baseURL string
}

type ctxKey string

const userIDKey ctxKey = "user_id"

func NewClientHandler(s client.ClientService, log logger.Logger, baseURL string) *ClientHandler {
	return &ClientHandler{
		service: s,
		log:     log,
		baseURL: baseURL,
	}
}

// @Summary List clients
// @Description Получить список клиентов с фильтрацией и пагинацией
// @Tags clients
// @Produce json
// @Param status query string false "Client status"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} dto.ListClientsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clients [get]
func (h *ClientHandler) List(c *gin.Context) {
	var q dto.ListClientQuery

	h.log.Info("http list clients started")

	if err := c.ShouldBindQuery(&q); err != nil {
		h.log.Warn("invalid query params")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query params",
		})
		return
	}

	if q.Limit == 0 {
		q.Limit = 20
	}

	list, total, err := h.service.List(
		c.Request.Context(),
		q.Status,
		q.Limit,
		q.Offset,
	)

	if err != nil {
		h.log.Error("list clients failed",
			"status", q.Status,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list clients",
		})
		return
	}

	h.log.Info("clients listed",
		"count", len(list),
		"total", total,
		"status", q.Status,
	)

	c.JSON(http.StatusOK, dto.ListClientsResponse{
		Items: mapper.ToClientSummaryList(list),
		Total: total,
	})
}

// @Summary Create API client
// @Description Создание нового API клиента
// @Tags clients
// @Accept json
// @Produce json
// @Param request body dto.CreateClientRequest true "Client data"
// @Success 201 {object} dto.ClientResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clients [post]
func (h *ClientHandler) Create(c *gin.Context) {
	var req dto.CreateClientRequest
	userID := c.GetString("user_id")

	h.log.Info("http create client started",
		"user_id", userID,
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid create client body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	client, err := h.service.Create(
		c.Request.Context(),
		userID,
		req.Name,
		req.APIServiceID,
		req.Description,
	)

	if err != nil {
		h.log.Error("create client failed",
			"user_id", userID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create client",
		})
		return
	}

	h.log.Info("client created",
		"client_id", client.ID,
		"user_id", userID,
	)

	url := fmt.Sprintf(
		"%s/api/%s",
		h.baseURL,
		client.Slug,
	)

	resp := dto.ClientResponse{
		ID:           client.ID,
		Name:         client.Name,
		Slug:         client.Slug,
		URL:          url,
		APIServiceID: client.APIServiceID,
		Status:       string(client.GetStatus()),
		CreatedAt:    client.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, resp)
}

// @Summary Get client by ID
// @Description Получить клиента по ID
// @Tags clients
// @Produce json
// @Param client_id path string true "Client ID"
// @Success 200 {object} dto.ClientResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clients/{client_id} [get]
func (h *ClientHandler) GetByID(c *gin.Context) {
	clientID := c.Param("client_id")

	h.log.Info("http get client",
		"client_id", clientID,
	)

	if clientID == "" {
		h.log.Warn("missing client_id")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id is required",
		})
		return
	}

	client, err := h.service.GetByID(
		c.Request.Context(),
		clientID,
	)

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

		h.log.Error("get client failed",
			"client_id", clientID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch client",
		})
		return
	}

	url := fmt.Sprintf(
		"%s/api/%s",
		h.baseURL,
		client.Slug,
	)

	resp := dto.ClientResponse{
		ID:           client.ID,
		Name:         client.Name,
		Slug:         client.Slug,
		URL:          url,
		APIServiceID: client.APIServiceID,
		Status:       string(client.GetStatus()),
		CreatedAt:    client.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Restart client
// @Description Перезапуск клиента через обновление Deployment (rolling update)
// @Tags clients
// @Produce json
// @Param client_id path string true "Client ID"
// @Success 202 {object} dto.ClientActionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clients/{client_id}/restart [post]
func (h *ClientHandler) RestartById(c *gin.Context) {
	clientID := c.Param("client_id")
	userID := c.GetString("user_id")

	h.log.Info("http restart client started",
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

	err := h.service.Restart(
		c.Request.Context(),
		userID,
		clientID,
		string(domain.ActionRestart),
	)

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

		if errors.Is(err, domain.ErrInvalidStateTransition) {
			h.log.Warn("invalid state transition",
				"client_id", clientID,
			)

			c.JSON(http.StatusConflict, gin.H{"error": "invalid state transition"})
			return
		}

		h.log.Error("restart failed",
			"client_id", clientID,
			"user_id", userID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to restart client",
		})
		return
	}

	h.log.Info("restart scheduled",
		"client_id", clientID,
		"user_id", userID,
	)

	resp := dto.ClientActionResponse{
		ClientID: clientID,
		Action:   string(domain.ActionRestart),
		Status:   string(domain.ClientStatusRestarting),
	}

	c.JSON(http.StatusAccepted, resp)
}

// @Summary Delete client
// @Description Удаление клиента
// @Tags clients
// @Produce json
// @Param client_id path string true "Client ID"
// @Success 202 {object} dto.ClientActionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clients/{client_id}/delete [post]
func (h *ClientHandler) DeleteById(c *gin.Context) {
	clientID := c.Param("client_id")
	userID := c.GetString("user_id")

	h.log.Info("http delete client started",
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

	err := h.service.Delete(
		c.Request.Context(),
		userID,
		clientID,
	)

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

		h.log.Error("delete client failed",
			"client_id", clientID,
			"user_id", userID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete client",
		})
		return
	}

	h.log.Info("client delete scheduled",
		"client_id", clientID,
		"user_id", userID,
	)

	resp := dto.ClientActionResponse{
		ClientID: clientID,
		Action:   string(domain.ActionDelete),
		Status:   string(domain.ClientStatusDeleting),
	}

	c.JSON(http.StatusAccepted, resp)
}

// @Summary Start client
// @Description Запуск клиента
// @Tags clients
// @Produce json
// @Param client_id path string true "Client ID"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /clients/{client_id}/start [post]
func (h *ClientHandler) StartByID(c *gin.Context) {
	clientID := c.Param("client_id")
	userID := c.GetString("user_id")

	h.log.Info("http start client started",
		"client_id", clientID,
		"user_id", userID,
	)

	ctx := context.WithValue(c.Request.Context(), userIDKey, userID)

	if err := h.service.Start(ctx, userID, clientID); err != nil {
		h.log.Warn("start client failed",
			"client_id", clientID,
			"user_id", userID,
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.log.Info("client start initiated",
		"client_id", clientID,
		"user_id", userID,
	)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "starting",
	})
}

// @Summary Stop client
// @Description Остановка клиента
// @Tags clients
// @Produce json
// @Param client_id path string true "Client ID"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /clients/{client_id}/stop [post]
func (h *ClientHandler) StopByID(c *gin.Context) {
	clientID := c.Param("client_id")
	userID := c.GetString("user_id")

	h.log.Info("http stop client started",
		"client_id", clientID,
		"user_id", userID,
	)

	ctx := context.WithValue(c.Request.Context(), userIDKey, userID)

	if err := h.service.Stop(ctx, userID, clientID); err != nil {
		h.log.Warn("stop client failed",
			"client_id", clientID,
			"user_id", userID,
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.log.Info("client stop initiated",
		"client_id", clientID,
		"user_id", userID,
	)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "stopping",
	})
}
