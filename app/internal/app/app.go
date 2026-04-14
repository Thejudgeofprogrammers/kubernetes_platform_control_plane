package app

import (
	"context"
	"control_plane/internal/config"
	reconciler "control_plane/internal/reconciler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewApp(env *config.Config) (*gin.Engine, reconciler.ReconcilerService) {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://172.22.4.66:5173",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	rec := RegisterRoutes(r, env)
	ctx := context.Background()
	rec.Start(ctx)
	return r, rec
}
