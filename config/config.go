package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/pkg/errors"
)

type Config struct {
	Server *Server `envPrefix:"SERVER_"`
}

type Server struct {
	Port         string `env:"PORT"  envDefault:""`
	BaseUrl      string `env:"BASE_URL" envDefault:"127.0.0.1"`
	AppUrlPrefix string `env:"SERVER_ADDRESS" envDefault:"api/v1"`
}

func NewConfig() *Config {
	return &Config{
		Server: &Server{},
	}
}

func LoadConfig() (*Config, error) {
	cfg := NewConfig()
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	return cfg, nil
}
