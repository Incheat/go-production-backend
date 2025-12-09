// Package main defines the main function for the auth service.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	servergen "github.com/incheat/go-playground/services/auth/internal/api/gen/oapi/public/server"
	"github.com/incheat/go-playground/services/auth/internal/config"
	"github.com/incheat/go-playground/services/auth/internal/controller/auth"
	handler "github.com/incheat/go-playground/services/auth/internal/handler/http"
	memoryrepo "github.com/incheat/go-playground/services/auth/internal/repository/memory"
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
	swagger, err := servergen.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("failed to connect to redis: %w", err))
	}
	logger.Info("Redis connected", zap.String("addr", cfg.Redis.Addr))
	defer func() {
		if err := rdb.Close(); err != nil {
			logger.Info("failed to close redis: %v", zap.String("error", err.Error()))
		}
	}()

	r := chi.NewRouter()
	r.Use(nethttpmiddleware.OapiRequestValidator(swagger))

	refreshTokenRepo := memoryrepo.NewRefreshTokenRepository()
	jwt := token.NewJWTMaker(cfg.JWT.Secret, cfg.JWT.Expire)
	opaque := token.NewOpaqueMaker(cfg.Refresh.NumBytes, cfg.Refresh.MaxAge, cfg.Refresh.EndPoint)
	ctrl := auth.NewController(jwt, opaque, refreshTokenRepo)
	srv := handler.NewHandler(ctrl)
	strictServer := servergen.NewStrictHandler(srv, nil)
	h := servergen.HandlerFromMux(strictServer, r)

	var g errgroup.Group

	g.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), h)
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

// func convertCORSRules(cfg *config.Config) []globalmiddleware.CORSRule {
// 	corsRules := make([]globalmiddleware.CORSRule, len(cfg.CORS.Rules))
// 	for i, rule := range cfg.CORS.Rules {
// 		corsRules[i] = globalmiddleware.CORSRule{
// 			Path:           rule.Path,
// 			AllowedOrigins: rule.AllowedOrigins,
// 		}
// 	}
// 	return corsRules
// }
