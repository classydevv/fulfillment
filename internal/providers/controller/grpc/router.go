package grpc

import (
	v1 "github.com/classydevv/fulfillment/internal/providers/controller/grpc/v1"
	"github.com/classydevv/fulfillment/internal/providers/usecase"
	"github.com/classydevv/fulfillment/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouterProvider(s *grpc.Server, uc usecase.Provider, l logger.Interface) {
	{
		v1.NewControllerProvider(s, uc, l)
	}

	reflection.Register(s)
}
