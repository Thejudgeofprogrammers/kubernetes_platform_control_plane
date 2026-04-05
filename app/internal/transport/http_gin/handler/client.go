package handler

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/service/client"
	dto "control_plane/internal/transport/http_gin/dto/client"
	"control_plane/internal/transport/http_gin/mapper"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	service client.ClientService
	log     *slog.Logger
}

func NewClientHandler(s client.ClientService, log *slog.Logger) *ClientHandler {
	return &ClientHandler{
		service: s,
		log:     log,
	}
}

func getUserID(c *gin.Context) string {
	userID := c.GetString("user_id")
	return userID
}

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

	resp := dto.ClientResponse{
		ID:           client.ID,
		Name:         client.Name,
		APIServiceID: client.APIServiceID,
		Status:       string(client.GetStatus()),
		CreatedAt:    client.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, resp)
}

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

	resp := dto.ClientResponse{
		ID:           client.ID,
		Name:         client.Name,
		APIServiceID: client.APIServiceID,
		Status:       string(client.GetStatus()),
		CreatedAt:    client.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}

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

	var req dto.RestartClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid restart body",
			"client_id", clientID,
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	err := h.service.Restart(
		c.Request.Context(),
		userID,
		clientID,
		req.Reason,
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
			"error": "failed to fetch restart",
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

func (h *ClientHandler) StartByID(c *gin.Context) {
	clientID := c.Param("client_id")
	userID := c.GetString("user_id")

	h.log.Info("http start client started",
		"client_id", clientID,
		"user_id", userID,
	)

	ctx := context.WithValue(c.Request.Context(), "userID", userID)

	if err := h.service.Start(ctx, clientID); err != nil {
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
