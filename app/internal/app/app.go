package app

import (
	"context"
	"control_plane/internal/app/bootstrap"
	"control_plane/internal/app/routes"
	"control_plane/internal/config"
	reconciler "control_plane/internal/reconciler"
	"control_plane/internal/transport/http_gin/middleware"
	"github.com/gin-gonic/gin"
)

func NewApp(env *config.Config) (*gin.Engine, reconciler.ReconcilerService) {
	r := gin.New()

	r.Use(middleware.CORSMiddleware())

	rec := Run(r, env)
	ctx := context.Background() // Потом поменять
	go rec.Run(ctx)

	return r, rec
}

func Run(r *gin.Engine, env *config.Config) reconciler.ReconcilerService {
	container := bootstrap.BuildContainer(env)

	routes.RegisterRoutes(r, container, env)

	return container.Reconciler
}
