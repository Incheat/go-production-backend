package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/incheat/go-playground/internal/middleware"
	servergen "github.com/incheat/go-playground/services/helloworld/internal/api/gen/server"
	"github.com/incheat/go-playground/services/helloworld/internal/config"
	"github.com/incheat/go-playground/services/helloworld/internal/handler"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("ENV: %s\n", os.Getenv("APP_ENV"))
	fmt.Printf("Server port: %d\n", cfg.Server.Port)
	fmt.Printf("DB host: %s\n", cfg.Database.Host)

	r := gin.Default()
	r.Use(middleware.PathBasedCORS(convertCORSRules(cfg)))

	srv := handler.NewServer()
	handler := servergen.NewStrictHandler(srv, nil)
	servergen.RegisterHandlers(r, handler)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
	}

	log.Fatal(s.ListenAndServe())
	
}

func convertCORSRules(cfg *config.Config) []middleware.CORSRule {
	corsRules := make([]middleware.CORSRule, len(cfg.CORS.Rules))
	for i, rule := range cfg.CORS.Rules {
		corsRules[i] = middleware.CORSRule{
			Path:           rule.Path,
			AllowedOrigins: rule.AllowedOrigins,
		}
	}
	return corsRules
}