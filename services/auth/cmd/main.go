package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	globalmiddleware "github.com/incheat/go-playground/internal/middleware"
	"github.com/incheat/go-playground/internal/oapi"
	servergen "github.com/incheat/go-playground/services/auth/internal/api/gen/server"
	"github.com/incheat/go-playground/services/auth/internal/config"
	"github.com/incheat/go-playground/services/auth/internal/controller/auth"
	httphandler "github.com/incheat/go-playground/services/auth/internal/handler/http"
	localmiddleware "github.com/incheat/go-playground/services/auth/internal/middleware"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("ENV: %s\n", cfg.Env)
	fmt.Printf("Server port: %d\n", cfg.Server.Port)

	logger, _ := zap.NewDevelopment()

	// Get OpenAPI definition from embedded spec
	swagger, err := servergen.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec: %v", err)
	}

	switch cfg.Env {
	case config.EnvDev, config.EnvStaging:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

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

	ctrl := auth.NewController(nil, nil, nil)
	srv := httphandler.NewHandler(ctrl)
	handler := servergen.NewStrictHandler(srv, nil)
	servergen.RegisterHandlers(r, handler)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
	}

	log.Fatal(s.ListenAndServe())
	
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