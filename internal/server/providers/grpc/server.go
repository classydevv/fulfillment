package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	config "github.com/classydevv/fulfillment/configs/providers"
	pb "github.com/classydevv/fulfillment/pkg/api/providers"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type providerId string

type Server struct {
	pb.UnimplementedProvidersServiceServer

	grpc struct {
		lis    net.Listener
		server *grpc.Server
	}

	mu        sync.RWMutex
	providers map[providerId]*pb.Provider
}

// New - returns *Server
func New(ctx context.Context, cfg *config.Config) (*Server, error) {
	srv := &Server{
		providers: make(map[providerId]*pb.Provider),
	}

	// grpc
	{
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
		if err != nil {
			return nil, fmt.Errorf("server failed to listen: %w", err)
		}

		s := grpc.NewServer()
		pb.RegisterProvidersServiceServer(s, srv)

		reflection.Register(s)

		srv.grpc.lis = lis
		srv.grpc.server = s
	}

	// grpc gateway
	{

	}

	// swagger-ui
	{

	}

	return srv, nil
}

func (s *Server) Run(ctx context.Context) error {
	group := errgroup.Group{}
	group.Go(func() error {
		log.Printf("start serve at: %v", s.grpc.lis.Addr())
		if err := s.grpc.server.Serve(s.grpc.lis); err != nil {
			return fmt.Errorf("failed to serve server: %w", err)
		}
		return nil
	})

	return group.Wait()
}
