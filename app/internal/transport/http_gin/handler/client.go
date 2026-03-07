package handler

import (
	"control_plane/internal/domain"
	"control_plane/internal/service/client"
	dto "control_plane/internal/transport/http_gin/dto/client"
	"control_plane/internal/transport/http_gin/mapper"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	service client.ClientService
}

func NewClientHandler(s client.ClientService) *ClientHandler {
	return &ClientHandler{
		service: s,
	}
}

func getUserID(c *gin.Context) string {
	userID := c.GetString("user_id")
	return userID
}

func (h *ClientHandler) List(c *gin.Context) {
	var q dto.ListClientQuery

	if err := c.ShouldBindQuery(&q); err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list clients",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ListClientsResponse{
		Items: mapper.ToClientSummaryList(list),
		Total: total,
	})
}

func (h *ClientHandler) Create(c *gin.Context) {
	var req dto.CreateClientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")

	client, err := h.service.Create(
		c.Request.Context(),
		userID,
		req.Name,
		req.APIServiceID,
		req.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create client",
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

	c.JSON(http.StatusCreated, resp)
}

func (h *ClientHandler) GetByID(c *gin.Context) {
	clientID := c.Param("client_id")

	if clientID == "" {
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
			c.JSON(http.StatusNotFound, gin.H{
				"error": "client not found",
			})
			return
		}
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

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id is required",
		})
		return
	}

	var req dto.RestartClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	userID := c.GetString("user_id")

	err := h.service.Restart(
		c.Request.Context(),
		userID,
		clientID,
		req.Reason,
	)

	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "client not found",
			})
			return
		}

		if errors.Is(err, domain.ErrInvalidStateTransition) {
			c.JSON(http.StatusConflict, gin.H{"error": "invalid state transition"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch restart",
		})
		return
	}

	resp := dto.ClientActionResponse{
		ClientID: clientID,
		Action:   string(domain.ActionRestart),
		Status:   string(domain.ClientStatusRestarting),
	}

	c.JSON(http.StatusAccepted, resp)
}

func (h *ClientHandler) DeleteById(c *gin.Context) {
	clientID := c.Param("client_id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client id is required",
		})
		return
	}

	userID := c.GetString("user_id")

	err := h.service.Delete(
		c.Request.Context(),
		userID,
		clientID,
	)

	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "client not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete client",
		})
		return
	}

	resp := dto.ClientActionResponse{
		ClientID: clientID,
		Action:   string(domain.ActionDelete),
		Status:   string(domain.ClientStatusDeleting),
	}

	c.JSON(http.StatusAccepted, resp)
}
