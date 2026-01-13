// Package main defines the main function for the auth service.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	servergen "github.com/incheat/go-production-backend/services/auth/internal/api/oapi/gen/public/server"
	envconfig "github.com/incheat/go-production-backend/services/auth/internal/config/env"
	usergateway "github.com/incheat/go-production-backend/services/auth/internal/gateway/user/grpc"
	authhandler "github.com/incheat/go-production-backend/services/auth/internal/handler/http"
	chimiddleware "github.com/incheat/go-production-backend/services/auth/internal/middleware/chi"
	redisrepo "github.com/incheat/go-production-backend/services/auth/internal/repository/redis"
	authservice "github.com/incheat/go-production-backend/services/auth/internal/service/auth"
	"github.com/incheat/go-production-backend/services/auth/internal/token"
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
		jwksPath = "/.well-known/jwks.json"
	}
	rootRouter.Get(jwksPath, jwtTokenMaker.JWKSHandler)

	// HTTP API router
	apiRouter := chi.NewRouter()
	apiRouter.Use(nethttpmiddleware.OapiRequestValidatorWithOptions(
		openAPISpec,
		chimiddleware.NewValidatorOptions(chimiddleware.ValidatorConfig{
			ProdMode: cfg.Env == envconfig.EnvProd,
		}),
	))
	// apiRouter.Use(chimiddleware.PathBasedCORS(convertCORSRules(&cfg.CORS.Internal)))
	// apiRouter.Use(chimiddleware.RequestID())
	apiRouter.Use(chimiddleware.HTTPRequest())
	apiRouter.Use(chimiddleware.ZapLogger(logger))
	apiRouter.Use(chimiddleware.ZapRecovery(logger))

	apiHandler := servergen.HandlerFromMux(strict, apiRouter)

	rootRouter.Mount("/", apiHandler)

	var g errgroup.Group

	g.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%d", int(cfg.Server.PublicPort)), rootRouter)
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

// func convertCORSRules(corsRule *envconfig.CORSRule) []chimiddleware.CORSRule {
// 	path := "*"
// 	return []chimiddleware.CORSRule{
// 		{
// 			Path:           path,
// 			AllowedOrigins: corsRule.AllowedOrigins,
// 		},
// 	}
// }
