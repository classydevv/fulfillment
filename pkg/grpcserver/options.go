package grpcserver

import (
	"net"
)

type Option func(*Server)

func AddressGRPC(host, port string) Option {
	return func(s *Server) {
		s.GRPC.Address = net.JoinHostPort(host, port)
	}
}

func AddressGateway(host, port string) Option {
	return func(s *Server) {
		s.Gateway.Address = net.JoinHostPort(host, port)
	}
}
