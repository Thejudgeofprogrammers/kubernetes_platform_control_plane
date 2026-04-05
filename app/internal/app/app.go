package app

import (
	"control_plane/internal/config"
	reconciler "control_plane/internal/reconciler"

	"github.com/gin-gonic/gin"
)

func NewApp(env *config.Config) (*gin.Engine, reconciler.ReconcilerService) {
	r := gin.New()

	rec := RegisterRoutes(r, env)

	return r, rec
}
