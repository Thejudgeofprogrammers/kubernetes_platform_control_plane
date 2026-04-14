package handler

import (
	"control_plane/internal/service/user"
	authDTO "control_plane/internal/transport/http_gin/dto/auth"
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

// @Summary List users
// @Description Получить список пользователей
// @Tags users
// @Produce json
// @Success 200 {object} map[string][]domain.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
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

// @Summary Delete user
// @Description Удаление пользователя
// @Tags users
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{user_id} [delete]
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
			"error": "delete user failed",
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

// @Summary Update user role
// @Description Обновление роли пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body dto.UpdateRoleRequest true "Role data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /users/{user_id}/role [patch]
func (h *UserHandler) UpdateRole(c *gin.Context) {
	userID := c.Param("user_id")
	requestUserID := c.GetString("user_id")

	h.log.Info("http update role started",
		"target_user_id", userID,
		"requested_by", requestUserID,
	)

	var req authDTO.UpdateRoleRequest

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
			"error": "update role failed",
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

// @Summary Get current user
// @Description Получить текущего пользователя
// @Tags users
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} map[string]string
// @Router /users/me [get]
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
			"error": "get me failed",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
