package httpserver

import (
	"net"
	"time"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

func Prefork(prefork bool) Option {
	return func(s *Server) {
		s.prefork = prefork
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.httpReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.httpWriteTimeout = timeout
	}
}

func ServerShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.serverShutdownTimeout = timeout
	}
}
