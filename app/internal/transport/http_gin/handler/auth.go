package handler

import (
	"control_plane/internal/service/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service auth.AuthService
}

func NewAuthHandler(s auth.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		FullName string `json:"full_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	err := h.service.Register(
		c.Request.Context(),
		req.Email,
		req.FullName,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "registered",
	})
}

func (h *AuthHandler) RequestCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}

	err := h.service.RequestCode(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to request code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "code_sent",
	})
}

func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid code",
		})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
