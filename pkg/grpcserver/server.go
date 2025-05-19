package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	GRPC struct {
		Address string
		Server  *grpc.Server
	}
	Gateway struct {
		Address string
		Server  *http.Server
		Mux     *runtime.ServeMux
	}
	notify chan error
}

func New(ctx context.Context, opts ...Option) *Server {
	s := &Server{
		notify: make(chan error, 10),
	}

	s.GRPC.Server = grpc.NewServer()
	mux := runtime.NewServeMux()
	s.Gateway.Mux = mux
	s.Gateway.Server = &http.Server{Handler: mux}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run() {
	go func() {
		lis, err := net.Listen("tcp", s.GRPC.Address)
		if err != nil {
			s.notify <- fmt.Errorf("grpc - New - Run - GRPC - net.Listen: %w", err)
			return
		}
		s.notify <- s.GRPC.Server.Serve(lis)
	}()
	go func() {
		lis, err := net.Listen("tcp", s.Gateway.Address)
		if err != nil {
			s.notify <- fmt.Errorf("grpc - New - Run - Gateway - net.Listen: %w", err)
			return
		}
		s.notify <- s.Gateway.Server.Serve(lis)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.GRPC.Server.GracefulStop()
	s.Gateway.Server.Shutdown(ctx)

	return nil
}
