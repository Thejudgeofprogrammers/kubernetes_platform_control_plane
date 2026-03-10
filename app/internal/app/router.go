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
	health "control_plane/internal/service/health/impl"
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
	codeRepository := memory.NewInMemoryEmailCodeRepository()
	healthRepository := memory.NewInMemoryClientHealthRepository()
	// apiServiceRepository := memory.NewInMemoryAPIServiceRepository()

	// redis
	rdb := redis.NewRedisClient(env.RedisAddr, env.GetRedisPassword(), env.RedisDB)

	// Services
	auditService := audit.NewAuditService(actionRepository)
	healthService := health.NewHealthService(healthRepository)
	orchestrator := mock.NewMockOrchestrator(healthService)
	configService := cfgService.NewConfigService(
		clientRepository,
		configRepository,
		auditService,
		orchestrator,
	)
	clientService := client.NewClientService(clientRepository, orchestrator, auditService)
	jwtService := jwt.NewJWTService(env.GetSecret(), env.Exp)
	refreshService := refresh.NewRefreshService(rdb, env.Ref_time)
	authService := auth.NewAuthService(userRepository, codeRepository, refreshService, jwtService, env.ExpireEmailCode)

	// Handlers
	clientHandler := handler.NewClientHandler(clientService)
	configHandler := handler.NewConfigHandler(configService)
	authHandler := handler.NewAuthHandler(authService)
	healthHandler := handler.NewHealthHandler(healthService)

	api := r.Group("/api/" + env.VersionAPI)
	protected := api.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(jwtService))

	// auth
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/request-code", authHandler.RequestCode)
	api.POST("/auth/verify-code", authHandler.VerifyCode)
	api.POST("/auth/refresh", authHandler.Refresh)

	// health
	protected.GET("/clients/:client_id/health", healthHandler.GetHealth)

	// client
	protected.GET("/clients", clientHandler.List)
	protected.GET("/clients/:client_id", clientHandler.GetByID)
	protected.GET("/clients/:client_id/configs", configHandler.ListConfigs)

	protected.POST(
		"/clients",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		clientHandler.Create,
	)
	protected.POST(
		"/clients/:client_id/restart",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		clientHandler.RestartById,
	)
	protected.POST(
		"/clients/:client_id/delete",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		clientHandler.DeleteById,
	)
	protected.POST(
		"/clients/:client_id/configs",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		configHandler.CreateConfig,
	)
	protected.POST(
		"/clients/:client_id/configs/:config_id/deploy",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		configHandler.Deploy,
	)
}
