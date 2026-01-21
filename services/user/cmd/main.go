// Package main defines the main function for the user service.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	userpb "github.com/incheat/go-production-backend/api/user/grpc/gen"
	envconfig "github.com/incheat/go-production-backend/services/user/internal/config/env"
	"github.com/incheat/go-production-backend/services/user/internal/constant"
	userhandler "github.com/incheat/go-production-backend/services/user/internal/handler/grpc"
	"github.com/incheat/go-production-backend/services/user/internal/interceptor"
	"github.com/incheat/go-production-backend/services/user/internal/obs"
	userrepo "github.com/incheat/go-production-backend/services/user/internal/repository/mysql"
	userservice "github.com/incheat/go-production-backend/services/user/internal/service/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {

	cfg, err := envconfig.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	logger := initLogger(envconfig.EnvName(cfg.Env))

	logger.Info("Starting user service", zap.String("env", string(cfg.Env)))
	logger.Info("GRPC server internal port", zap.Int("port", int(cfg.Server.InternalPort)))

	// ------------------------------------------------------------------
	// Context & signal
	// ------------------------------------------------------------------
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ------------------------------------------------------------------
	// OpenTelemetry (MUST be before creating grpcServer)
	// ------------------------------------------------------------------
	otelShutdown, err := obs.InitTracer(ctx, constant.ServiceName, cfg.OTEL.Endpoint)
	if err != nil {
		log.Fatalf("Error initializing OpenTelemetry tracer: %v", err)
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			logger.Error("Error shutting down OpenTelemetry tracer", zap.Error(err))
		}
	}()

	interceptors := interceptor.DefaultChain(logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithFilter(filters.Not(filters.HealthCheck())),
			),
		),
	)

	// ----------------------------
	// gRPC Health Service
	// ----------------------------
	healthServer := health.NewServer()

	// Initially not serving, wait for MySQL to be ready
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	// Initialize MySQL connection
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.DBName)
	logger.Info("Initializing MySQL connection", zap.String("dsn", dbDSN))
	dbConn, err := sql.Open("mysql", dbDSN)
	if err != nil {
		log.Fatalf("Error opening MySQL connection: %v", err)
	}
	dbConn.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	dbConn.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	dbConn.SetConnMaxLifetime(time.Duration(cfg.MySQL.ConnMaxLifetime) * time.Second)
	{
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := dbConn.PingContext(ctx); err != nil {
			healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			log.Fatalf("Error pinging MySQL: %v", err)
		}
	}

	// ✅ MySQL OK -> readiness OK
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	defer func() {
		// Before closing, declare NOT_SERVING to stop traffic from Envoy/K8s
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

		grpcServer.GracefulStop()

		if err := dbConn.Close(); err != nil {
			logger.Warn("Failed to close MySQL connection", zap.Error(err))
		}
	}()

	// user components
	userRepository := userrepo.NewUserRepository(dbConn)

	userService := userservice.New(userRepository)
	userImpl := userhandler.New(userService)

	userpb.RegisterUserServiceInternalServer(grpcServer, userImpl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.InternalPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var g errgroup.Group

	g.Go(func() error {
		return grpcServer.Serve(lis)
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
