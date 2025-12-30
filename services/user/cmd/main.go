// Package main defines the main function for the user service.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	userpb "github.com/incheat/go-production-backend/api/user/grpc/gen"
	envconfig "github.com/incheat/go-production-backend/services/user/internal/config/env"
	userhandler "github.com/incheat/go-production-backend/services/user/internal/handler/grpc"
	"github.com/incheat/go-production-backend/services/user/internal/interceptor"
	userrepo "github.com/incheat/go-production-backend/services/user/internal/repository/mysql"
	userservice "github.com/incheat/go-production-backend/services/user/internal/service/user"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := envconfig.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	logger := initLogger(envconfig.EnvName(cfg.Env))

	logger.Info("Starting user service", zap.String("env", string(cfg.Env)))
	logger.Info("GRPC server internal port", zap.Int("port", int(cfg.Server.InternalPort)))

	limiter := rate.NewLimiter(100, 100)
	interceptors := interceptor.DefaultChain(limiter)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)

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
