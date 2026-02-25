package app

import (
	"control_plane/internal/config"
	"control_plane/internal/domain"
	"control_plane/internal/orchestrator/mock"
	"control_plane/internal/repository/memory"
	"control_plane/internal/service"
	"control_plane/internal/transport/http_gin/handler"
	"control_plane/internal/transport/http_gin/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, env *config.Config) {
	// Repositories
	clientRepository := memory.NewInMemoryClientRepository()
	actionRepository := memory.NewInMemoryClientActionRepository()
	configRepository := memory.NewInMemoryClientConfigRepository()

	// mock
	orchestrator := mock.NewMockOrchestrator()

	// Services
	auditService := service.NewAuditService(actionRepository)
	configService := service.NewConfigService(
		clientRepository,
		configRepository,
		auditService,
		orchestrator,
	)

	clientService := service.NewClientService(clientRepository, orchestrator, auditService)

	// Handlers ............
	clientHandler := handler.NewClientHandler(clientService)
	configHandler := handler.NewConfigHandler(configService)

	api := r.Group("/api/v1")

	api.Use(
		middleware.AddUserID(),
		middleware.AuthMiddleware(env.AllowUnauthorized),
	)

	api.GET("/clients", clientHandler.List)
	api.GET("/clients/:client_id", clientHandler.GetByID)
	api.GET("/clients/:client_id/configs", configHandler.ListConfigs)

	api.POST(
		"/clients",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		clientHandler.Create,
	)
	api.POST(
		"/clients/:client_id/restart",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		clientHandler.RestartById,
	)
	api.POST(
		"/clients/:client_id/delete",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		clientHandler.DeleteById,
	)
	api.POST(
		"/clients/:client_id/configs",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		configHandler.CreateConfig,
	)
	api.POST(
		"/clients/:client_id/configs/:config_id/deploy",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		configHandler.Deploy,
	)
}
