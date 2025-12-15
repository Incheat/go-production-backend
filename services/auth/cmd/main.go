// Package main defines the main function for the auth service.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	globalchimiddleware "github.com/incheat/go-playground/internal/middleware/chi"
	servergen "github.com/incheat/go-playground/services/auth/internal/api/oapi/gen/public/server"
	envconfig "github.com/incheat/go-playground/services/auth/internal/config/env"
	usergateway "github.com/incheat/go-playground/services/auth/internal/gateway/user/http"
	authhandler "github.com/incheat/go-playground/services/auth/internal/handler/http"
	chimiddleware "github.com/incheat/go-playground/services/auth/internal/middleware/chi"
	redisrepo "github.com/incheat/go-playground/services/auth/internal/repository/redis"
	authservice "github.com/incheat/go-playground/services/auth/internal/service/auth"
	"github.com/incheat/go-playground/services/auth/internal/token"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {

	cfg, err := envconfig.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logger := initLogger(cfg.Env)

	logger.Info("Starting auth service", zap.String("env", string(cfg.Env)))
	logger.Info("Http server port", zap.Int("port", int(cfg.Server.PublicPort)))

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

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("failed to connect to Redis: %w", err))
	}

	logger.Info("Connected to Redis", zap.String("addr", cfg.Redis.Host))
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Warn("Failed to close Redis client", zap.Error(err))
		}
	}()

	// HTTP router
	router := chi.NewRouter()
	router.Use(nethttpmiddleware.OapiRequestValidatorWithOptions(
		openAPISpec,
		globalchimiddleware.NewValidatorOptions(globalchimiddleware.ValidatorConfig{
			ProdMode: cfg.Env == envconfig.EnvProd,
		}),
	))
	router.Use(chimiddleware.PathBasedCORS(convertCORSRules(&cfg.CORS.Internal)))
	router.Use(chimiddleware.RequestID())
	router.Use(chimiddleware.HTTPRequest())
	router.Use(chimiddleware.ZapLogger(logger))
	router.Use(chimiddleware.ZapRecovery(logger))

	// Auth components
	refreshTokenRepository := redisrepo.NewRefreshTokenRepository(redisClient)

	jwtTokenMaker := token.NewJWTMaker(cfg.JWT.Secret, cfg.JWT.Expire)
	opaqueTokenMaker := token.NewOpaqueMaker(
		cfg.Refresh.NumBytes,
		cfg.Refresh.MaxAge,
		cfg.Refresh.EndPoint,
	)

	userGateway, err := usergateway.NewUserGateway(fmt.Sprintf("http://localhost:%d", cfg.UserGateway.InternalPort))
	if err != nil {
		log.Fatalf("Error creating user gateway: %v", err)
	}
	authService := authservice.New(jwtTokenMaker, opaqueTokenMaker, refreshTokenRepository, userGateway)
	authImpl := authhandler.New(authService)

	strict := servergen.NewStrictHandler(authImpl, nil)
	apiHandler := servergen.HandlerFromMux(strict, router)

	var g errgroup.Group

	g.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%d", int(cfg.Server.PublicPort)), apiHandler)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}

func initLogger(env envconfig.EnvName) *zap.Logger {
	switch env {
	case envconfig.EnvDev, envconfig.EnvStaging:
		return zap.Must(zap.NewDevelopment())
	default:
		return zap.Must(zap.NewProduction())
	}
}

func convertCORSRules(corsRule *envconfig.CORSRule) []chimiddleware.CORSRule {
	path := "*"
	return []chimiddleware.CORSRule{
		{
			Path:           path,
			AllowedOrigins: corsRule.AllowedOrigins,
		},
	}
}
