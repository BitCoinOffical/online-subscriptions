package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Env string

const (
	EnvProd = "Prod"
	EnvDev  = "Dev"
)

type Config struct {
	Postgres PostgresConfig
	App      AppConfig
}

type PostgresConfig struct {
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
}

type AppConfig struct {
	DebugLevel string `env:"DEBUG_LEVEL"`
}

func NewLoadConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg.Postgres); err != nil {
		return nil, err
	}
	if err := env.Parse(&cfg.App); err != nil {
		return nil, err
	}

	var env Env = Env(cfg.App.DebugLevel)
	if env != EnvProd && env != EnvDev {
		return nil, fmt.Errorf("incorrect debug level: %s", env)
	}
	return &cfg, nil
}
