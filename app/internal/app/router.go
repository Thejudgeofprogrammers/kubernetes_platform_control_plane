package app

import (
	"context"
	"control_plane/internal/app/bootstrap"
	"control_plane/internal/config"
	"control_plane/internal/domain"
	"control_plane/internal/infra/redis"
	logger "control_plane/internal/logger"
	"control_plane/internal/orchestrator"
	"control_plane/internal/orchestrator/k8s"
	reconciler "control_plane/internal/reconciler"
	reconciler_impl "control_plane/internal/reconciler/impl"
	"control_plane/internal/repository/memory"
	action "control_plane/internal/service/action/impl"
	audit "control_plane/internal/service/audit/impl"
	auth "control_plane/internal/service/auth/impl"
	client "control_plane/internal/service/client/impl"
	cfgService "control_plane/internal/service/config/impl"
	"control_plane/internal/service/email"
	emailImpl "control_plane/internal/service/email/impl"
	emailMock "control_plane/internal/service/email/mock"
	health "control_plane/internal/service/health/impl"
	jwt "control_plane/internal/service/jwt/impl"
	refresh "control_plane/internal/service/refresh/impl"
	user "control_plane/internal/service/user/impl"
	apiservice "control_plane/internal/service/api_service/impl"
	"time"

	"control_plane/internal/transport/http_gin/handler"
	"control_plane/internal/transport/http_gin/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, env *config.Config) reconciler.ReconcilerService {
	log := logger.New()

	// Repositories
	apiServiceRepository := memory.NewInMemoryAPIServiceRepository(log)
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
	apiServiceService := apiservice.NewAPIServiceService(apiServiceRepository, log)
	healthService := health.NewHealthService(healthRepository, log)

	clientset, err := k8s.NewClient()
	var orchestrator orchestrator.Orchestrator
	if err != nil {
		panic("K8s not started")
	} else {
		orchestrator = k8s.NewK8sOrchestrator(
			clientset, 
			env.Namespace, 
			env.ImageAPIClient, 
			env.ProxyConnectTimeout,
			env.ProxyReadTimeout,
			env.ProxySendTimeout,
			log, healthService, 
			apiServiceService,
		)
		log.Info("Запуск orchestrator в K8s")
	}

	var emailSender email.EmailSender

	if env.EnvFile == "prod" {
		emailSender = emailImpl.NewSMTPEmailSender(
			env.SMTPHost,
			env.SMTPPort,
			env.SMTPUser,
			env.SMTPPass,
			env.SMTPFrom,
		)
	} else {
		emailSender = emailMock.NewEmailSenderMock(log)
	}

	rec := reconciler_impl.NewReconciler(actionRepository, clientRepository, apiServiceRepository, orchestrator, configRepository, log, healthService)
	configService := cfgService.NewConfigService(clientRepository, configRepository, auditService, orchestrator, log)
	actionService := action.NewActionService(actionRepository, log)
	clientService := client.NewClientService(clientRepository, orchestrator, auditService, configRepository, actionService, log)
	jwtService := jwt.NewJWTService(env.GetSecret(), env.Exp, log)
	refreshService := refresh.NewRefreshService(rdb, env.Ref_time, log)
	authService := auth.NewAuthService(userRepository, codeRepository, refreshService, jwtService, emailSender, env.ExpireEmailCode, log)
	userService := user.NewUserService(userRepository, log)

	// #### Надо доделать prod dev, когда будет бд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bootstrap.SeedAdmin(ctx, userRepository, env.EmailAdmin, env.FullNameAdmin)

	// Handlers
	clientHandler := handler.NewClientHandler(clientService, log)
	configHandler := handler.NewConfigHandler(configService, log)
	authHandler := handler.NewAuthHandler(authService, log)
	healthHandler := handler.NewHealthHandler(healthService, log)
	userHandler := handler.NewUserHandler(userService, log)
	apiServiceHandler := handler.NewAPIServiceHandler(apiServiceService, log)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API
	api := r.Group("/api/" + env.VersionAPI)
	protected := api.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(jwtService))

	// auth
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/request-code", authHandler.RequestCode)
	api.POST("/auth/verify-code", authHandler.VerifyCode)
	api.POST("/auth/refresh", authHandler.Refresh)

	//metrics
	api.POST("/clients/:id/metrics")

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

	// api-service
	protected.GET("/api-services", apiServiceHandler.List)

	protected.POST(
		"/api-services",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		apiServiceHandler.Create,
	)

	protected.DELETE(
		"/api-services/:id",
		middleware.RoleMiddleware(env.AllowForbidden, string(domain.RoleOwner)),
		apiServiceHandler.Delete,
	)

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
