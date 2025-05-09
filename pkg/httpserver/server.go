package httpserver

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App    *fiber.App
	notify chan error

	address               string
	prefork               bool
	httpReadTimeout       time.Duration
	httpWriteTimeout      time.Duration
	serverShutdownTimeout time.Duration
}

func New(opts ...Option) *Server {
	s := &Server{
		App:    nil,
		notify: make(chan error, 1),
	}

	for _, opt := range opts {
		opt(s)
	}

	app := fiber.New(fiber.Config{
		Prefork:      s.prefork,
		ReadTimeout:  s.httpReadTimeout,
		WriteTimeout: s.httpWriteTimeout,
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
	return s.App.ShutdownWithTimeout(s.serverShutdownTimeout)
}
