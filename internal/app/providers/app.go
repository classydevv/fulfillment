package providers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config "github.com/classydevv/fulfillment/configs/providers"
	router "github.com/classydevv/fulfillment/internal/controller/http"
	"github.com/classydevv/fulfillment/internal/repo/persistent"
	grpcserver "github.com/classydevv/fulfillment/internal/server/providers/grpc"
	httpserver "github.com/classydevv/fulfillment/internal/server/providers/http"
	"github.com/classydevv/fulfillment/internal/usecase/provider"
	"github.com/classydevv/fulfillment/pkg/logger"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// ** Repository **

	// ** UseCase **
	providerUseCase := provider.New(
		persistent.New(),
	)

	// ** Delivery **
	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
	router.NewProviderRouter(httpServer.App, cfg, l, providerUseCase)
	
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
