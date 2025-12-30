// Package main defines the main function for the helloworld service.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	servergen "github.com/incheat/go-production-backend/services/helloworld/internal/api/oapi/gen/public/server"
	koanfconfig "github.com/incheat/go-production-backend/services/helloworld/internal/config/koanf"
	"github.com/incheat/go-production-backend/services/helloworld/internal/handler"
	middleware "github.com/incheat/go-production-backend/services/helloworld/internal/middleware/gin"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
)

func main() {
	cfg := koanfconfig.MustLoad()
	env := os.Getenv("APP_ENV")

	fmt.Printf("ENV: %s\n", env)
	fmt.Printf("Server port: %d\n", cfg.Server.Port)
	fmt.Printf("DB host: %s\n", cfg.Database.Host)

	// Get OpenAPI definition from embedded spec
	swagger, err := servergen.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec: %v", err)
	}

	r := gin.Default()
	// Apply CORS rules based on the request path.
	r.Use(middleware.PathBasedCORS(convertCORSRules(cfg)))
	// Validate requests against the OpenAPI schema.
	r.Use(ginmiddleware.OapiRequestValidatorWithOptions(
		swagger,
		middleware.NewValidatorOptions(middleware.ValidatorConfig{
			ProdMode: env == "prod",
		}),
	))

	srv := handler.NewServer()
	handler := servergen.NewStrictHandler(srv, nil)
	servergen.RegisterHandlers(r, handler)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
	}

	log.Fatal(s.ListenAndServe())

}

func convertCORSRules(cfg *koanfconfig.Config) []middleware.CORSRule {
	corsRules := make([]middleware.CORSRule, len(cfg.CORS.Rules))
	for i, rule := range cfg.CORS.Rules {
		corsRules[i] = middleware.CORSRule{
			Path:           rule.Path,
			AllowedOrigins: rule.AllowedOrigins,
		}
	}
	return corsRules
}
