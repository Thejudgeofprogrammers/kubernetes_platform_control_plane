package handler

import (
	"control_plane/internal/service/auth"
	authDTO "control_plane/internal/transport/http_gin/dto/auth"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service auth.AuthService
	log     *slog.Logger
}

func NewAuthHandler(s auth.AuthService, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		service: s,
		log:     log,
	}
}

// @Summary Register user
// @Description Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req authDTO.RegisterRequest

	h.log.Info("http register started")

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid register request body")

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	err := h.service.Register(
		c.Request.Context(),
		req.Email,
		req.FullName,
	)

	if err != nil {
		h.log.Error("register failed",
			"email", req.Email,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register",
		})
		return
	}

	h.log.Info("register completed",
		"email", req.Email,
	)

	c.JSON(http.StatusCreated, gin.H{
		"status": "registered",
	})
}

// @Summary Request verification code
// @Description Отправка кода на email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RequestCodeRequest true "Email"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/request-code [post]
func (h *AuthHandler) RequestCode(c *gin.Context) {
	var req authDTO.RequestCodeRequest

	h.log.Info("http request code started")

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid request code body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}

	err := h.service.RequestCode(c.Request.Context(), req.Email)
	if err != nil {
		h.log.Error("request code failed",
			"email", req.Email,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to request code",
		})
		return
	}

	h.log.Info("code requested",
		"email", req.Email,
	)

	c.JSON(http.StatusOK, gin.H{
		"status": "code_sent",
	})
}

// @Summary Verify code
// @Description Подтверждение кода и получение JWT токенов
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyCodeRequest true "Verification data"
// @Success 200 {object} domain.AuthTokens
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/verify-code [post]
func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var req authDTO.VerifyCodeRequest

	h.log.Info("http verify code started")

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid verify code body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}

	tokens, err := h.service.VerifyCode(
		c.Request.Context(),
		req.Email,
		req.Code,
	)

	if err != nil {
		h.log.Warn("verify code failed",
			"email", req.Email,
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid code",
		})
		return
	}

	h.log.Info("verify code success",
		"email", req.Email,
	)

	c.JSON(http.StatusOK, authDTO.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary Refresh token
// @Description Обновление access токена
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh token"
// @Success 200 {object} domain.AuthTokens
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req authDTO.RefreshRequest

	h.log.Info("http refresh started")

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid refresh body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}

	tokens, err := h.service.Refresh(
		c.Request.Context(),
		req.RefreshToken,
	)

	if err != nil {
		h.log.Warn("refresh failed")

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	h.log.Info("refresh success")

	c.JSON(http.StatusOK, authDTO.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
