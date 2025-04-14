package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		GRPC GRPC
		HTTP HTTP
		Log Log
		Metrics Metrics
		Swagger Swagger
	}

	GRPC struct {
		Port string `env:"GRPC_PORT,required"`
	}

	HTTP struct {
		// Host string `env:"HTTP_HOST,required"`
		Port string `env:"HTTP_PORT,required"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"debug"`
	}

	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"true"`
	}

)

func NewConfig() (*Config, error) {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}