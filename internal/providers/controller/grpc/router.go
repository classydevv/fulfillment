package grpc

import (
	"context"

	v1 "github.com/classydevv/fulfillment/internal/providers/controller/grpc/v1"
	"github.com/classydevv/fulfillment/internal/providers/usecase"
	"github.com/classydevv/fulfillment/pkg/grpcserver"
	"github.com/classydevv/fulfillment/pkg/logger"
	"google.golang.org/grpc/reflection"
)

func NewRouterProvider(ctx context.Context, s *grpcserver.Server, uc usecase.Provider, l logger.Interface) {
	{
		v1.NewControllerProvider(ctx, s, uc, l)
	}

	reflection.Register(s.GRPC.Server)
}
