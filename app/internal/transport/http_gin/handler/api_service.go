package handler

import (
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	apiservice "control_plane/internal/service/api_service"
	apiserviceDTO "control_plane/internal/transport/http_gin/dto/api_service"
	"control_plane/internal/transport/http_gin/mapper"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIServiceHandler struct {
	service apiservice.APIServiceService
	log     logger.Logger
}

func NewAPIServiceHandler(s apiservice.APIServiceService, log logger.Logger) *APIServiceHandler {
	return &APIServiceHandler{
		service: s,
		log:     log,
	}
}

// @Summary List API services
// @Tags api-services
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.APIServiceResponse
// @Failure 500 {object} map[string]string
// @Router /api-services [get]
func (h *APIServiceHandler) List(c *gin.Context) {
	list, err := h.service.List(c.Request.Context())
	if err != nil {
		h.log.Error("list api services failed", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list api services",
		})
		return
	}

	resp := make([]apiserviceDTO.APIServiceResponse, 0, len(list))
	for _, s := range list {
		resp = append(resp, mapper.ToAPIServiceResponse(s))
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Create API service
// @Tags api-services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAPIServiceRequest true "API service data"
// @Success 201 {object} dto.APIServiceResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api-services [post]
func (h *APIServiceHandler) Create(c *gin.Context) {
	var req apiserviceDTO.CreateAPIServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	service, err := h.service.Create(
		c.Request.Context(),
		req.Name,
		req.BaseURL,
		req.Protocol,
	)
	if err != nil {
		h.log.Error("create api service failed", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create api service",
		})
		return
	}

	c.JSON(http.StatusCreated, mapper.ToAPIServiceResponse(service))
}

// @Summary Delete API service
// @Tags api-services
// @Produce json
// @Security BearerAuth
// @Param id path string true "Service ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api-services/{id} [delete]
func (h *APIServiceHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete api service",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

func (h *APIServiceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	service, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrAPIServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get api service",
		})
		return
	}

	c.JSON(http.StatusOK, mapper.ToAPIServiceResponse(service))
}

func (h *APIServiceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req apiserviceDTO.UpdateAPIServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	service, err := h.service.Update(
		c.Request.Context(),
		id,
		req.Name,
		req.BaseURL,
		req.Protocol,
	)
	if err != nil {
		if errors.Is(err, domain.ErrAPIServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update api service",
		})
		return
	}

	c.JSON(http.StatusOK, mapper.ToAPIServiceResponse(service))
}
