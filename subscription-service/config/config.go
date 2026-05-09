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
	Redis    RedisConfig
}

type PostgresConfig struct {
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBName     string `env:"DB_NAME,required"`
}
type RedisConfig struct {
	RDBAddr      string `env:"RDB_ADDR,required"`
	RDBPort      string `env:"RDB_PORT,required"`
	RDBPass      string `env:"RDB_PASS,required"`
	RDBLimiterDB int    `env:"RDB_RATE_LIMITER_DB,required"`
}
type AppConfig struct {
	Port       string `env:"PORT,required"`
	DebugLevel string `env:"DEBUG_LEVEL,required"`
	Jwt        string `env:"JWT_KEY,required"`
}

func NewLoadConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg.Postgres); err != nil {
		return nil, err
	}
	if err := env.Parse(&cfg.App); err != nil {
		return nil, err
	}
	if err := env.Parse(&cfg.Redis); err != nil {
		return nil, err
	}
	var env Env = Env(cfg.App.DebugLevel)
	if env != EnvProd && env != EnvDev {
		return nil, fmt.Errorf("incorrect debug level: %s", env)
	}
	return &cfg, nil
}
