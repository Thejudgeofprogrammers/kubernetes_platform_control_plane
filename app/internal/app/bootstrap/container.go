package bootstrap

import (
	"context"
	"control_plane/internal/config"
	"control_plane/internal/infra/redis"
	logger "control_plane/internal/logger"
	"control_plane/internal/orchestrator"
	"control_plane/internal/orchestrator/k8s"
	reconciler "control_plane/internal/reconciler"
	reconciler_impl "control_plane/internal/reconciler/impl"
	"control_plane/internal/repository/memory"
	action "control_plane/internal/service/action/impl"
	apiservice "control_plane/internal/service/api_service/impl"
	auth "control_plane/internal/service/auth/impl"
	client "control_plane/internal/service/client/impl"
	cfgService "control_plane/internal/service/config/impl"
	emailImpl "control_plane/internal/service/email/impl"
	health "control_plane/internal/service/health/impl"
	jwtSrv "control_plane/internal/service/jwt"
	jwt "control_plane/internal/service/jwt/impl"
	metric "control_plane/internal/service/metric/impl"
	refresh "control_plane/internal/service/refresh/impl"
	user "control_plane/internal/service/user/impl"
	"time"

	"control_plane/internal/transport/http_gin/handler"
)

type Container struct {
	ClientHandler     *handler.ClientHandler
	ConfigHandler     *handler.ConfigHandler
	AuthHandler       *handler.AuthHandler
	UserHandler       *handler.UserHandler
	HealthHandler     *handler.HealthHandler
	APIServiceHandler *handler.APIServiceHandler

	Reconciler reconciler.ReconcilerService
	JWTService jwtSrv.JWTService
}

func BuildContainer(env *config.Config) *Container {
	baseLog := logger.NewLogger(env.LogLevel, env.LogLayers)

	handlerLog := baseLog.With("layer", "handler")
	serviceLog := baseLog.With("layer", "service")
	repoLog := baseLog.With("layer", "repository")
	reconcilerLog := baseLog.With("layer", "reconciler")
	orchestratorLog := baseLog.With("layer", "orchestrator")

	// Repositories
	apiServiceRepository := memory.NewInMemoryAPIServiceRepository(repoLog)
	// clientAccessRepository := memory.NewInMemoryClientAccessRepository()

	clientRepository := memory.NewInMemoryClientRepository(repoLog)
	actionRepository := memory.NewInMemoryClientActionRepository(repoLog)
	configRepository := memory.NewInMemoryClientConfigRepository(repoLog)
	userRepository := memory.NewInMemoryUserRepository(repoLog)
	codeRepository := memory.NewInMemoryEmailCodeRepository(repoLog)
	healthRepository := memory.NewInMemoryClientHealthRepository(repoLog)
	metricRepository := memory.NewMetricsInMemory(env.MaxMetricsPerClient)

	// redis
	rdb := redis.NewRedisClient(env.RedisAddr, env.GetRedisPassword(), env.RedisDB)

	// Services
	actionService := action.NewActionService(actionRepository, serviceLog)
	apiServiceService := apiservice.NewAPIServiceService(apiServiceRepository, serviceLog)
	healthService := health.NewHealthService(healthRepository, serviceLog)
	metricService := metric.NewMetricsService()

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
			orchestratorLog,
			healthService,
			apiServiceService,
			clientRepository,
		)
		orchestratorLog.Info("Запуск orchestrator в K8s")
	}

	emailSender := emailImpl.NewSMTPEmailSender(
		env.SMTPHost,
		env.SMTPPort,
		env.SMTPUser,
		env.SMTPPass,
		env.SMTPFrom,
	)

	rec := reconciler_impl.NewReconciler(
		actionRepository,
		clientRepository,
		apiServiceRepository,
		orchestrator,
		configRepository,
		reconcilerLog,
		healthService,
		metricRepository,
		metricService,
		env.BaseURLIngress,
	)
	configService := cfgService.NewConfigService(clientRepository, configRepository, actionService, orchestrator, serviceLog)
	clientService := client.NewClientService(clientRepository, orchestrator, configRepository, actionService, serviceLog)
	jwtService := jwt.NewJWTService(env.GetSecret(), env.Exp, serviceLog)
	refreshService := refresh.NewRefreshService(rdb, env.Ref_time, serviceLog)
	authService := auth.NewAuthService(userRepository, codeRepository, refreshService, jwtService, emailSender, env.ExpireEmailCode, serviceLog)
	userService := user.NewUserService(userRepository, serviceLog)

	// #### Надо доделать prod dev, когда будет бд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	SeedAdmin(ctx, userRepository, env.EmailAdmin, env.FullNameAdmin)

	// Handlers
	clientHandler := handler.NewClientHandler(clientService, handlerLog, env.PublicBaseURL)
	configHandler := handler.NewConfigHandler(configService, handlerLog)
	authHandler := handler.NewAuthHandler(authService, handlerLog)
	healthHandler := handler.NewHealthHandler(healthService, handlerLog)
	userHandler := handler.NewUserHandler(userService, handlerLog)
	apiServiceHandler := handler.NewAPIServiceHandler(apiServiceService, handlerLog)

	return &Container{
		ClientHandler:     clientHandler,
		ConfigHandler:     configHandler,
		AuthHandler:       authHandler,
		UserHandler:       userHandler,
		HealthHandler:     healthHandler,
		APIServiceHandler: apiServiceHandler,
		Reconciler:        rec,
		JWTService:        jwtService,
	}
}
