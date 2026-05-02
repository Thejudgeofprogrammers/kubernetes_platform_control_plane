package routes

import (
	"control_plane/internal/app/bootstrap"
	"control_plane/internal/config"
	"control_plane/internal/domain"
	"control_plane/internal/transport/http_gin/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, c *bootstrap.Container, env *config.Config) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/" + env.VersionAPI)
	protected := api.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(c.JWTService))

	api.POST("/auth/register", c.AuthHandler.Register)
	api.POST("/auth/request-code", c.AuthHandler.RequestCode)
	api.POST("/auth/verify-code", c.AuthHandler.VerifyCode)
	api.POST("/auth/refresh", c.AuthHandler.Refresh)
	api.POST("/clients/:client_id/metrics")

	protected.GET("/users", c.UserHandler.List)
	protected.GET("/users/me", c.UserHandler.Me)
	protected.GET("/clients/:client_id/health", c.HealthHandler.GetHealth)
	protected.GET("/clients", c.ClientHandler.List)
	protected.GET("/clients/:client_id", c.ClientHandler.GetByID)
	protected.GET("/clients/:client_id/configs", c.ConfigHandler.ListConfigs)
	protected.GET("/api-services", c.APIServiceHandler.List)
	protected.GET("/api-services/:id", c.APIServiceHandler.GetByID)

	protected.DELETE(
		"/users/:user_id",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.UserHandler.Delete,
	)
	protected.PATCH(
		"/users/:user_id/role",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.UserHandler.UpdateRole,
	)

	protected.POST(
		"/api-services",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.APIServiceHandler.Create,
	)

	protected.PUT(
		"/api-services/:id",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.APIServiceHandler.Update,
	)

	protected.DELETE(
		"/api-services/:id",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.APIServiceHandler.Delete,
	)

	protected.POST(
		"/clients",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.ClientHandler.Create,
	)

	protected.POST(
		"/clients/:client_id/restart",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.ClientHandler.RestartById,
	)

	protected.POST(
		"/clients/:client_id/delete",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.ClientHandler.DeleteById,
	)

	protected.POST(
		"/clients/:client_id/configs",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.ConfigHandler.CreateConfig,
	)

	protected.POST(
		"/clients/:client_id/configs/:config_id/deploy",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.ConfigHandler.Deploy,
	)

	protected.PUT(
		"/clients/:client_id/configs/:config_id/update",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.ConfigHandler.Update,
	)

	protected.DELETE(
		"/clients/:client_id/configs/:config_id/delete",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.ConfigHandler.Delete,
	)

	protected.POST(
		"/clients/:client_id/start",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.ClientHandler.StartByID,
	)

	protected.POST(
		"/clients/:client_id/stop",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		c.ClientHandler.StopByID,
	)
}
