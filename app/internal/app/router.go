package app

import (
	"control_plane/internal/config"
	"control_plane/internal/domain"
	"control_plane/internal/infra/redis"
	"control_plane/internal/orchestrator/mock"
	"control_plane/internal/repository/memory"
	audit "control_plane/internal/service/audit/impl"
	auth "control_plane/internal/service/auth/impl"
	client "control_plane/internal/service/client/impl"
	cfgService "control_plane/internal/service/config/impl"
	jwt "control_plane/internal/service/jwt/impl"
	refresh "control_plane/internal/service/refresh/impl"
	"control_plane/internal/transport/http_gin/handler"
	"control_plane/internal/transport/http_gin/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, env *config.Config) {
	// Repositories
	clientRepository := memory.NewInMemoryClientRepository()
	actionRepository := memory.NewInMemoryClientActionRepository()
	configRepository := memory.NewInMemoryClientConfigRepository()
	userRepository := memory.NewInMemoryUserRepository()
	// emailCodeRepository := memory.NewInMemoryEmailCodeRepository()

	// mock
	orchestrator := mock.NewMockOrchestrator()

	// redis
	rdb := redis.NewRedisClient(env.RedisAddr, env.GetRedisPassword(), env.RedisDB)

	// Services
	auditService := audit.NewAuditService(actionRepository)
	configService := cfgService.NewConfigService(
		clientRepository,
		configRepository,
		auditService,
		orchestrator,
	)

	clientService := client.NewClientService(clientRepository, orchestrator, auditService)
	jwtService := jwt.NewJWTService(env.GetSecret(), env.Exp)
	refreshService := refresh.NewRefreshService(rdb)
	authService := auth.NewAuthService(userRepository, refreshService, jwtService)

	// Handlers ............
	clientHandler := handler.NewClientHandler(clientService)
	configHandler := handler.NewConfigHandler(configService)

	api := r.Group("/api/v1")

	api.Use(
		middleware.AddUserID(),
		// middleware.JWTAuthMiddleware(),
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
