// Package main defines the main function for the auth service.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	obsconfig "github.com/incheat/go-production-backend/pkg/obs/config"
	"github.com/incheat/go-production-backend/pkg/obs/logging"
	obsmetrics "github.com/incheat/go-production-backend/pkg/obs/metrics"
	"github.com/incheat/go-production-backend/pkg/obs/profiling"
	obstracing "github.com/incheat/go-production-backend/pkg/obs/tracing"
	servergen "github.com/incheat/go-production-backend/services/auth/internal/api/oapi/gen/public/server"
	envconfig "github.com/incheat/go-production-backend/services/auth/internal/config/env"
	"github.com/incheat/go-production-backend/services/auth/internal/constant"
	usergateway "github.com/incheat/go-production-backend/services/auth/internal/gateway/user/grpc"
	authhandler "github.com/incheat/go-production-backend/services/auth/internal/handler/http"
	chimiddleware "github.com/incheat/go-production-backend/services/auth/internal/middleware/chi"
	redisrepo "github.com/incheat/go-production-backend/services/auth/internal/repository/redis"
	authservice "github.com/incheat/go-production-backend/services/auth/internal/service/auth"
	"github.com/incheat/go-production-backend/services/auth/internal/token"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {

	cfg, err := envconfig.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	telemetryConfig := obsconfig.TelemetryConfig{
		Resource: obsconfig.ResourceConfig{
			ServiceName:    constant.ServiceName,
			Environment:    string(cfg.Env),
			ServiceVersion: cfg.Version,
		},
		Logging: obsconfig.LoggingConfig{
			Level: cfg.Obs.Logging.Level,
		},
		OTLP: obsconfig.OTLPConfig{
			Endpoint: cfg.Obs.OTLP.Endpoint,
			Insecure: true,
		},
		Tracing: obsconfig.TracingConfig{
			SamplingRatio: cfg.Obs.Tracing.SamplingRatio,
		},
	}

	logger, err := logging.New(logging.Config{
		Service: telemetryConfig.Resource.ServiceName,
		Env:     telemetryConfig.Resource.Environment,
		Level:   telemetryConfig.Logging.Level,
	})
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	logger.Info("Starting auth service", zap.String("env", string(cfg.Env)))
	logger.Info("Http server port", zap.Int("port", int(cfg.Server.HTTPPort)))

	ctx := context.Background()

	// Start profiling server
	profiling.StartServer(ctx, fmt.Sprintf(":%d", int(cfg.Obs.Profiling.Port)), logger)

	// Initialize OpenTelemetry tracer
	otelShutdown, err := obstracing.InitTracer(ctx, telemetryConfig)
	if err != nil {
		logger.Error("Error initializing OpenTelemetry tracer", zap.Error(err))
	} else {
		logger.Info("OpenTelemetry tracer initialized", zap.String("endpoint", cfg.Obs.OTLP.Endpoint))
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			logger.Error("Error shutting down OpenTelemetry tracer", zap.Error(err))
		}
	}()

	// Initialize Prometheus metrics
	reg := obsmetrics.NewRegistry()
	obsmetrics.RegisterHTTP(reg)
	shutdownMetrics := obsmetrics.StartServer(fmt.Sprintf(":%d", int(cfg.Obs.Metrics.Port)), reg, logger)
	if err != nil {
		logger.Error("Error initializing Prometheus metrics server", zap.Error(err))
	} else {
		logger.Info("Prometheus metrics server initialized", zap.String("port", fmt.Sprintf(":%d", int(cfg.Obs.Metrics.Port))))
	}
	defer func() {
		if err := shutdownMetrics(ctx); err != nil {
			logger.Error("Error shutting down Prometheus metrics server", zap.Error(err))
		}
	}()

	// Get OpenAPI definition from embedded spec
	openAPISpec, err := servergen.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading OpenAPI spec: %v", err)
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("failed to connect to Redis: %w", err))
	}

	logger.Info("Connected to Redis", zap.String("addr", cfg.Redis.Host))
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Warn("Failed to close Redis client", zap.Error(err))
		}
	}()

	// Auth components
	refreshTokenRepository := redisrepo.NewRefreshTokenRepository(redisClient)

	jwtTokenMaker, err := token.New(cfg.JWT.PrivateKeyPEM, cfg.JWT.KeyID, cfg.JWT.Issuer, cfg.JWT.Audience, cfg.JWT.Expire)
	if err != nil {
		log.Fatalf("Error creating JWT token maker: %v", err)
	}
	opaqueTokenMaker := token.NewOpaqueMaker(
		cfg.Refresh.NumBytes,
		cfg.Refresh.MaxAge,
		cfg.Refresh.EndPoint,
	)

	logger.Info("Creating user gateway", zap.String("address", cfg.UserGateway.InternalAddress))
	userGateway, err := usergateway.New(cfg.UserGateway.InternalAddress)
	if err != nil {
		log.Fatalf("Error creating user gateway: %v", err)
	}
	authService := authservice.New(jwtTokenMaker, opaqueTokenMaker, refreshTokenRepository, userGateway)
	authImpl := authhandler.New(authService)

	strict := servergen.NewStrictHandler(authImpl, nil)

	// ---- HTTP Routers ----
	rootRouter := chi.NewRouter()

	// ✅ Health check endpoint
	rootRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if _, err := w.Write([]byte("redis not ready")); err != nil {
				logger.Error("Failed to write health check response", zap.Error(err))
			}
			return
		}

		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("ok")); err != nil {
			logger.Error("Failed to write health check response", zap.Error(err))
		}
	})

	// ✅ JWKS endpoint (NOT behind OpenAPI validator)
	jwksPath := cfg.JWT.JWKSPath
	if jwksPath == "" {
		jwksPath = constant.JWKSPath
	}
	rootRouter.Get(jwksPath, jwtTokenMaker.JWKSHandler)

	// ✅ Traced router
	tracedRouter := chi.NewRouter()

	tracedRouter.Use(obsmetrics.PromHTTPMetrics())

	// HTTP API router
	apiRouter := chi.NewRouter()
	apiRouter.Use(chimiddleware.RequestMeta())
	apiRouter.Use(logging.HTTPRequestLogging(logger))
	apiRouter.Use(nethttpmiddleware.OapiRequestValidatorWithOptions(
		openAPISpec,
		chimiddleware.NewValidatorOptions(chimiddleware.ValidatorConfig{
			ProdMode: cfg.Env == envconfig.EnvProd,
		}),
	))
	apiRouter.Use(chimiddleware.ZapRecovery(logger))

	apiHandler := servergen.HandlerFromMux(strict, apiRouter)

	tracedRouter.Mount("/", apiHandler)

	rootRouter.Mount("/", otelhttp.NewHandler(
		tracedRouter,
		constant.SpanNameAuthHTTP,
	))

	var g errgroup.Group

	listenAddr := fmt.Sprintf(":%d", int(cfg.Server.HTTPPort))

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: rootRouter,
	}

	g.Go(func() error {
		return srv.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}

// func initLogger(env envconfig.EnvName) *zap.Logger {
// 	switch env {
// 	case envconfig.EnvDev, envconfig.EnvStaging:
// 		return zap.Must(zap.NewDevelopment())
// 	default:
// 		return zap.Must(zap.NewProduction())
// 	}
// }
