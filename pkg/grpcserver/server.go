package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	App  *grpc.Server
	notify chan error
	address string
}

func New(opts ...Option) *Server {
	s := &Server {
		App: grpc.NewServer(),
		notify: make(chan error, 1),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run() {
	go func() {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			s.notify <- fmt.Errorf("grpc - New - net.Listen: %w", err)
			close(s.notify)
		}
		s.notify <- s.App.Serve(lis)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	s.App.GracefulStop()

	return nil
}
