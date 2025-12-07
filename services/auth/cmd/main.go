// Package main defines the main function for the auth service.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	globalmiddleware "github.com/incheat/go-playground/internal/middleware"
	"github.com/incheat/go-playground/internal/oapi"
	servergen "github.com/incheat/go-playground/services/auth/internal/api/gen/oapi/public/server"
	"github.com/incheat/go-playground/services/auth/internal/config"
	"github.com/incheat/go-playground/services/auth/internal/controller/auth"
	handler "github.com/incheat/go-playground/services/auth/internal/handler/http"
	localmiddleware "github.com/incheat/go-playground/services/auth/internal/middleware"
	"github.com/incheat/go-playground/services/auth/internal/token"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
	"go.uber.org/zap"
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

	r := initGin(cfg.Env)

	// Apply CORS rules based on the request path.
	r.Use(
		globalmiddleware.PathBasedCORS(convertCORSRules(cfg)),
		localmiddleware.ZapLogger(logger),
		localmiddleware.ZapRecovery(logger),
		localmiddleware.RequestID(),
	)

	// Validate requests against the OpenAPI schema.
	r.Use(ginmiddleware.OapiRequestValidatorWithOptions(
		swagger,
		oapi.NewValidatorOptions(oapi.ValidatorConfig{
			ProdMode: cfg.Env == config.EnvProd,
		}),
	))

	jwt := token.NewJWTMaker(cfg.JWT.Secret, cfg.JWT.Expire)
	opaque := token.NewOpaqueMaker(cfg.Refresh.NumBytes, cfg.Refresh.MaxAge, cfg.Refresh.EndPoint)
	ctrl := auth.NewController(jwt, opaque, nil)
	srv := handler.NewHandler(ctrl)
	handler := servergen.NewStrictHandler(srv, nil)
	servergen.RegisterHandlers(r, handler)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
	}

	log.Fatal(s.ListenAndServe())
}

func initLogger(env config.EnvName) *zap.Logger {
	switch env {
	case config.EnvDev, config.EnvStaging:
		return zap.Must(zap.NewDevelopment())
	default:
		return zap.Must(zap.NewProduction())
	}
}

func initGin(env config.EnvName) *gin.Engine {
	switch env {
	case config.EnvDev, config.EnvStaging:
		gin.SetMode(gin.DebugMode)
		return gin.New()
	default:
		gin.SetMode(gin.ReleaseMode)
		return gin.New()
	}
}

func convertCORSRules(cfg *config.Config) []globalmiddleware.CORSRule {
	corsRules := make([]globalmiddleware.CORSRule, len(cfg.CORS.Rules))
	for i, rule := range cfg.CORS.Rules {
		corsRules[i] = globalmiddleware.CORSRule{
			Path:           rule.Path,
			AllowedOrigins: rule.AllowedOrigins,
		}
	}
	return corsRules
}
