package server

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	// "github.com/goccy/go-json"
)

const (
	_defaultAddr            = ":8080"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	App *fiber.App
	notify chan error

	address string
	prefork bool
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

func New(opts ...Option) *Server {
	s := &Server{
		App: nil,
		notify: make(chan error, 1),
		address:         _defaultAddr,
		readTimeout: _defaultReadTimeout,
		writeTimeout: _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	app := fiber.New(fiber.Config{
		Prefork: s.prefork,
		ReadTimeout: s.readTimeout,
		WriteTimeout: s.writeTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	s.App = app

	return s
}

func (s *Server) Run() {
	go func() {
		s.notify <- s.App.Listen(s.address)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.App.ShutdownWithTimeout(s.shutdownTimeout)
}