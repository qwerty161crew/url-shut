package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Config struct {
	Server   *Server   `envPrefix:"SERVER_"`
	Postgres *Postgres `envPrefix:"POSTGRES_"`
	Security *Security `envPrefix:"SECURITY_"`
}
type Security struct {
	Salt string `env:"SALT"  envDefault:"123e4567-e89b-"`
}
type Logger struct {
	LogLevel string `env:"LEVEL"  envDefault:""`
}
type Postgres struct {
	Host     string `env:"HOST"  envDefault:"127.0.0.1"`
	Port     string `env:"PORT"  envDefault:"5432"`
	User     string `env:"USER"  envDefault:"postgres"`
	Password string `env:"PASSWORD"  envDefault:"1917"`
	Db       string `env:"PASSWORD"  envDefault:"url"`
}

type Server struct {
	Mode         string `env:"MODE"      envDefault:"local"`
	LogFile      string `env:"LOG_FILE"   envDefault:"output.log"`
	Port         string `env:"PORT"  envDefault:""`
	BaseUrl      string `env:"BASE_URL" envDefault:"127.0.0.1"`
	AppUrlPrefix string `env:"SERVER_ADDRESS" envDefault:"/api/v1"`
	File         string `env:"URL_FILE" envDefault:"url.txt"`
}

func NewConfig() *Config {
	return &Config{
		Server:   &Server{},
		Postgres: &Postgres{},
		Security: &Security{},
	}
}

func LoadConfig() (*Config, error) {
	cfg := NewConfig()
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	return cfg, nil
}
func (s *Server) Level() zerolog.Level {
	if s.Mode != "prod" {
		return zerolog.DebugLevel
	}
	return zerolog.InfoLevel
}

func (s *Server) OutputFile() string {
	return s.LogFile
}

func (pg *Postgres) GenerateDBurl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pg.User,
		pg.Password,
		pg.Host,
		pg.Port,
		pg.Db,
	)
}
