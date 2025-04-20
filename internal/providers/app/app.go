package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/classydevv/fulfillment/configs/providers"
	httpcontroller "github.com/classydevv/fulfillment/internal/providers/controller/http"
	repo "github.com/classydevv/fulfillment/internal/providers/repo/persistent/postgres"
	"github.com/classydevv/fulfillment/internal/providers/usecase"
	grpcserver "github.com/classydevv/fulfillment/pkg/grpc"
	httpserver "github.com/classydevv/fulfillment/pkg/http"
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
		postgres.MaxPoolSize(int32(cfg.PG.MaxPoolSize)))
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
		httpserver.ReadTimeout(time.Duration(cfg.HTTP.ReadTimeoutSeconds) * time.Second),
		httpserver.WriteTimeout(time.Duration(cfg.HTTP.WriteTimeoutSeconds) * time.Second),
		httpserver.ServerShutdownTimeout(time.Duration(cfg.HTTP.ServerShutdownTimeout) * time.Second),
	)
	httpcontroller.NewRouterProvider(httpServer.App, cfg, l, providerUseCase)
	
	// GRPC Server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcServer, err := grpcserver.New(ctx, cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("providers - Run - grpcserver.New: %s", err))
	}

	// Start servers
	httpServer.Run()
	// redo with channel!!
	if err = grpcServer.Run(ctx); err != nil {
		l.Fatal(fmt.Errorf("providers - Run - grpcserver.Run: %w", err))
	}

	// Wait for errors
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("providers - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("providers - Run - httpServer.Notify: %w", err))
	}

	// Gracefull Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("providers - Run - httpServer.Shutdown: %w", err))
	}
}
