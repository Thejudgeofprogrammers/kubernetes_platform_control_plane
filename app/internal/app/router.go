package app

import (
	"control_plane/internal/config"
	"control_plane/internal/domain"
	"control_plane/internal/infra/redis"
	logger "control_plane/internal/logger"
	"control_plane/internal/orchestrator"
	"control_plane/internal/orchestrator/k8s"
	"control_plane/internal/orchestrator/mock"
	reconciler "control_plane/internal/reconciler"
	reconciler_impl "control_plane/internal/reconciler/impl"
	"control_plane/internal/repository/memory"
	action "control_plane/internal/service/action/impl"
	audit "control_plane/internal/service/audit/impl"
	auth "control_plane/internal/service/auth/impl"
	client "control_plane/internal/service/client/impl"
	cfgService "control_plane/internal/service/config/impl"
	health "control_plane/internal/service/health/impl"
	jwt "control_plane/internal/service/jwt/impl"
	refresh "control_plane/internal/service/refresh/impl"
	user "control_plane/internal/service/user/impl"
	"control_plane/internal/transport/http_gin/handler"
	"control_plane/internal/transport/http_gin/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, env *config.Config) reconciler.ReconcilerService {
	log := logger.New()

	// Repositories
	// apiServiceRepository := memory.NewInMemoryAPIServiceRepository()
	// clientAccessRepository := memory.NewInMemoryClientAccessRepository()

	clientRepository := memory.NewInMemoryClientRepository(log)
	actionRepository := memory.NewInMemoryClientActionRepository(log)
	configRepository := memory.NewInMemoryClientConfigRepository(log)
	userRepository := memory.NewInMemoryUserRepository(log)
	codeRepository := memory.NewInMemoryEmailCodeRepository(log)
	healthRepository := memory.NewInMemoryClientHealthRepository(log)

	// redis
	rdb := redis.NewRedisClient(env.RedisAddr, env.GetRedisPassword(), env.RedisDB)

	// Services
	auditService := audit.NewAuditService(actionRepository, log)
	healthService := health.NewHealthService(healthRepository, log)

	clientset, err := k8s.NewClient()
	var orchestrator orchestrator.Orchestrator
	if err != nil {
		// локально — mock
		orchestrator = mock.NewMockOrchestrator(healthService, log)
		log.Info("Запуск orchestrator, локально")
	} else {
		// в k8s
		orchestrator = k8s.NewK8sOrchestrator(clientset, env.Namespace, log)
		log.Info("Запуск orchestrator в K8s")
	}

	rec := reconciler_impl.NewReconciler(actionRepository, clientRepository, orchestrator, configRepository, log)
	configService := cfgService.NewConfigService(clientRepository, configRepository, auditService, orchestrator, log)
	actionService := action.NewActionService(actionRepository, log)
	clientService := client.NewClientService(clientRepository, orchestrator, auditService, configRepository, actionService, log)
	jwtService := jwt.NewJWTService(env.GetSecret(), env.Exp, log)
	refreshService := refresh.NewRefreshService(rdb, env.Ref_time, log)
	authService := auth.NewAuthService(userRepository, codeRepository, refreshService, jwtService, env.ExpireEmailCode, log)
	userService := user.NewUserService(userRepository, log)

	// Handlers
	clientHandler := handler.NewClientHandler(clientService, log)
	configHandler := handler.NewConfigHandler(configService, log)
	authHandler := handler.NewAuthHandler(authService, log)
	healthHandler := handler.NewHealthHandler(healthService, log)
	userHandler := handler.NewUserHandler(userService, log)

	// API
	api := r.Group("/api/" + env.VersionAPI)
	protected := api.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(jwtService))

	// auth
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/request-code", authHandler.RequestCode)
	api.POST("/auth/verify-code", authHandler.VerifyCode)
	api.POST("/auth/refresh", authHandler.Refresh)

	// user
	protected.GET("/users", userHandler.List)
	protected.DELETE("/users/:user_id", userHandler.Delete)
	protected.PATCH("/users/:user_id/role", userHandler.UpdateRole)
	protected.GET("/users/me", userHandler.Me)

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

	protected.POST(
		"/clients/:client_id/start",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		clientHandler.StartByID,
	)

	return rec
}
