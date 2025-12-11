// Package main defines the main function for the auth service.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	globalchimiddleware "github.com/incheat/go-playground/internal/middleware/chi"
	servergen "github.com/incheat/go-playground/services/auth/internal/api/gen/oapi/public/server"
	"github.com/incheat/go-playground/services/auth/internal/config"
	authhandler "github.com/incheat/go-playground/services/auth/internal/handler/http"
	customchimiddleware "github.com/incheat/go-playground/services/auth/internal/middleware/chi"
	redisrepo "github.com/incheat/go-playground/services/auth/internal/repository/redis"
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
	logger.Info("Connecting to Redis", zap.String("addr", cfg.Redis.Addr))

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

	// HTTP router
	router := chi.NewRouter()
	router.Use(nethttpmiddleware.OapiRequestValidatorWithOptions(
		openAPISpec,
		globalchimiddleware.NewValidatorOptions(globalchimiddleware.ValidatorConfig{
			ProdMode: cfg.Env == config.EnvProd,
		}),
	))
	router.Use(globalchimiddleware.PathBasedCORS(convertCORSRules(cfg)))
	router.Use(customchimiddleware.RequestID())
	router.Use(customchimiddleware.HTTPRequest())
	router.Use(customchimiddleware.ZapLogger(logger))
	router.Use(customchimiddleware.ZapRecovery(logger))
	router.Use(chimiddleware.Heartbeat("/healthz"))
	router.Use(chimiddleware.Heartbeat("/live"))
	router.Get("/ready", func(w http.ResponseWriter, _ *http.Request) {
		// check DB, cache, etc.
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			_, err := w.Write([]byte("not ready: " + err.Error()))
			if err != nil {
				logger.Error("failed to write response", zap.Error(err))
			}
			return
		}
		_, err := w.Write([]byte("ready"))
		if err != nil {
			logger.Error("failed to write response", zap.Error(err))
		}
	})

	// Auth components
	refreshTokenRepository := redisrepo.NewRefreshTokenRepository(redisClient)

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
