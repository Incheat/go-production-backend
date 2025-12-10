// Package main defines the main function for the auth service.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	globalchimiddleware "github.com/incheat/go-playground/internal/middleware/chi"
	servergen "github.com/incheat/go-playground/services/auth/internal/api/gen/oapi/public/server"
	"github.com/incheat/go-playground/services/auth/internal/config"
	authhandler "github.com/incheat/go-playground/services/auth/internal/handler/http"
	memoryrepo "github.com/incheat/go-playground/services/auth/internal/repository/memory"
	authservice "github.com/incheat/go-playground/services/auth/internal/service/auth"
	"github.com/incheat/go-playground/services/auth/internal/token"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {

	cfg := config.MustLoad()
	logger := initLogger(cfg.Env)

	logger.Info("Starting auth service", zap.String("env", string(cfg.Env)))
	logger.Info("Server port", zap.Int("port", cfg.Server.Port))

	// Get OpenAPI definition from embedded spec
	openAPISpec, err := servergen.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading OpenAPI spec: %v", err)
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("failed to connect to Redis: %w", err))
	}

	logger.Info("Connected to Redis", zap.String("addr", cfg.Redis.Addr))
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Warn("Failed to close Redis client", zap.Error(err))
		}
	}()

	// HTTP router
	router := chi.NewRouter()
	router.Use(nethttpmiddleware.OapiRequestValidator(openAPISpec))
	router.Use(globalchimiddleware.PathBasedCORS(convertCORSRules(cfg)))

	// Auth components
	refreshTokenRepository := memoryrepo.NewRefreshTokenRepository()

	jwtTokenMaker := token.NewJWTMaker(cfg.JWT.Secret, cfg.JWT.Expire)
	opaqueTokenMaker := token.NewOpaqueMaker(
		cfg.Refresh.NumBytes,
		cfg.Refresh.MaxAge,
		cfg.Refresh.EndPoint,
	)

	authService := authservice.New(jwtTokenMaker, opaqueTokenMaker, refreshTokenRepository)
	authImpl := authhandler.New(authService)

	strict := servergen.NewStrictHandler(authImpl, nil)
	apiHandler := servergen.HandlerFromMux(strict, router)

	var g errgroup.Group

	g.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), apiHandler)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}

func initLogger(env config.EnvName) *zap.Logger {
	switch env {
	case config.EnvDev, config.EnvStaging:
		return zap.Must(zap.NewDevelopment())
	default:
		return zap.Must(zap.NewProduction())
	}
}

func convertCORSRules(cfg *config.Config) []globalchimiddleware.CORSRule {
	corsRules := make([]globalchimiddleware.CORSRule, len(cfg.CORS.Rules))
	for i, rule := range cfg.CORS.Rules {
		corsRules[i] = globalchimiddleware.CORSRule{
			Path:           rule.Path,
			AllowedOrigins: rule.AllowedOrigins,
		}
	}
	return corsRules
}
