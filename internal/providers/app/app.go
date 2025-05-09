package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/classydevv/fulfillment/configs/providers"
	"github.com/classydevv/fulfillment/internal/providers/controller/http"
	"github.com/classydevv/fulfillment/internal/providers/controller/grpc"
	repo "github.com/classydevv/fulfillment/internal/providers/repo/persistent/postgres"
	"github.com/classydevv/fulfillment/internal/providers/usecase"
	"github.com/classydevv/fulfillment/pkg/grpcserver"
	"github.com/classydevv/fulfillment/pkg/httpserver"
	"github.com/classydevv/fulfillment/pkg/logger"
	"github.com/classydevv/fulfillment/pkg/postgres"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// ** Repository **
	pg, err := postgres.New(
		cfg.PG.URL,
		l,
		postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// ** UseCase **
	providerUseCase := usecase.NewUseCaseProviders(
		repo.NewPostgresRepo(pg),
	)

	// ** Delivery **
	// HTTP Server
	httpServer := httpserver.New(
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(time.Duration(cfg.HTTP.ReadTimeoutSeconds)*time.Second),
		httpserver.WriteTimeout(time.Duration(cfg.HTTP.WriteTimeoutSeconds)*time.Second),
		httpserver.ServerShutdownTimeout(time.Duration(cfg.HTTP.ServerShutdownTimeout)*time.Second),
	)
	http.NewRouterProvider(httpServer.App, providerUseCase, cfg, l)

	// GRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	grpc.NewRouterProvider(grpcServer.App, providerUseCase, l)

	// Start servers
	httpServer.Run()
	grpcServer.Run()

	// Wait for errors
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("providers - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("providers - Run - httpServer.Notify: %w", err))
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("providers - Run - grpcServer.Notify: %w", err))
	}

	// Graceful Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("providers - Run - httpServer.Shutdown: %w", err))
	}
	err = grpcServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("providers - Run - grpcServer.Shutdown: %w", err))
	}
}
