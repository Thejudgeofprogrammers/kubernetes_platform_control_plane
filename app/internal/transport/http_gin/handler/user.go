package handler

import (
	"control_plane/internal/service/user"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service user.UserService
	log     *slog.Logger
}

func NewUserHandler(s user.UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{
		service: s,
		log:     log,
	}
}

func (h *UserHandler) List(c *gin.Context) {
	h.log.Info("http list users started")

	users, err := h.service.List(c)
	if err != nil {
		h.log.Error("list users failed",
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list users",
		})
		return
	}

	h.log.Info("list users completed",
		"count", len(users),
	)

	c.JSON(http.StatusOK, gin.H{
		"items": users,
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	userID := c.Param("user_id")
	requestUserID := c.GetString("user_id")

	h.log.Info("http delete user started",
		"target_user_id", userID,
		"requested_by", requestUserID,
	)

	if err := h.service.Delete(c, userID); err != nil {
		h.log.Error("delete user failed",
			"target_user_id", userID,
			"requested_by", requestUserID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.log.Info("user deleted",
		"target_user_id", userID,
		"requested_by", requestUserID,
	)

	c.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	userID := c.Param("user_id")
	requestUserID := c.GetString("user_id")

	h.log.Info("http update role started",
		"target_user_id", userID,
		"requested_by", requestUserID,
	)

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid update role body",
			"target_user_id", userID,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}

	if err := h.service.UpdateRole(c, userID, req.Role); err != nil {
		h.log.Warn("update role failed",
			"target_user_id", userID,
			"requested_by", requestUserID,
			"role", req.Role,
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.log.Info("role updated",
		"target_user_id", userID,
		"requested_by", requestUserID,
		"role", req.Role,
	)

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")

	h.log.Debug("http get me",
		"user_id", userID,
	)

	user, err := h.service.GetMe(c, userID)
	if err != nil {
		h.log.Warn("get me failed",
			"user_id", userID,
			"error", err,
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
