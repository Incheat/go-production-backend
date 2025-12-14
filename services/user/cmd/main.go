// Package main defines the main function for the user service.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	globalchimiddleware "github.com/incheat/go-playground/internal/middleware/chi"
	servergen "github.com/incheat/go-playground/services/user/internal/api/oapi/gen/private/server"
	"github.com/incheat/go-playground/services/user/internal/config"
	userhandler "github.com/incheat/go-playground/services/user/internal/handler/http"
	chimiddleware "github.com/incheat/go-playground/services/user/internal/middleware/chi"
	userrepo "github.com/incheat/go-playground/services/user/internal/repository/mysql"
	userservice "github.com/incheat/go-playground/services/user/internal/service/user"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {

	cfg := config.MustLoad()
	logger := initLogger(cfg.Env)

	logger.Info("Starting user service", zap.String("env", string(cfg.Env)))
	logger.Info("Server port", zap.Int("port", cfg.Server.Port))

	// Get OpenAPI definition from embedded spec
	openAPISpec, err := servergen.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading OpenAPI spec: %v", err)
	}

	// HTTP router
	router := chi.NewRouter()
	router.Use(nethttpmiddleware.OapiRequestValidatorWithOptions(
		openAPISpec,
		globalchimiddleware.NewValidatorOptions(globalchimiddleware.ValidatorConfig{
			ProdMode: cfg.Env == config.EnvProd,
		}),
	))
	router.Use(globalchimiddleware.PathBasedCORS(convertCORSRules(cfg)))
	router.Use(chimiddleware.RequestID())
	router.Use(chimiddleware.ZapLogger(logger))
	router.Use(chimiddleware.ZapRecovery(logger))

	// Initialize MySQL connection
	logger.Info("Initializing MySQL connection", zap.String("dsn", cfg.MySQL.DSN))
	dbConn, err := sql.Open("mysql", cfg.MySQL.DSN)
	dbConn.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	dbConn.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	dbConn.SetConnMaxLifetime(time.Duration(cfg.MySQL.ConnMaxLifetime) * time.Second)
	if err != nil {
		log.Fatalf("Error opening MySQL connection: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			logger.Warn("Failed to close MySQL connection", zap.Error(err))
		}
	}()

	// Check if the connection is working
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Error pinging MySQL: %v", err)
	}

	// user components
	userRepository := userrepo.NewUserRepository(dbConn)

	userService := userservice.New(userRepository)
	userImpl := userhandler.New(userService)

	strict := servergen.NewStrictHandler(userImpl, nil)
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
