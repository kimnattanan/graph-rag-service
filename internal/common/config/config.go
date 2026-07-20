package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type CommonConfig struct {
	IsLocalEnv     bool   `env:"LOCAL_ENV" envDefault:"false"`
	Port           string `env:"HTTP_PORT,required"`
	AllowedOrigins string `env:"HTTP_ALLOWED_ORIGINS" envDefault:"http://localhost:3000"`
}

func NewConfig(cfg interface{}) error {
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("config error: %w", err)
	}
	return nil
}
