package config

import (
	"net"
	"net/url"
	"time"

	"github.com/kimnattanan/graph-rag-service/internal/common/config"
)

type (
	Config struct {
		Common   config.CommonConfig
		App      App
		Postgres Postgres
		Auth     Auth
		Admin    Admin
	}

	App struct {
		ServerToRun string `env:"SERVER_TO_RUN,required"`
	}

	Postgres struct {
		Host     string `env:"USER_POSTGRES_HOST" envDefault:"localhost"`
		Port     string `env:"USER_POSTGRES_PORT" envDefault:"5432"`
		User     string `env:"USER_POSTGRES_USER" envDefault:"postgres"`
		Password string `env:"USER_POSTGRES_PASSWORD" envDefault:"postgres"`
		Database string `env:"USER_POSTGRES_DB" envDefault:"users"`
		SSLMode  string `env:"USER_POSTGRES_SSLMODE" envDefault:"disable"`
	}

	Auth struct {
		JWTSecret       string `env:"JWT_SECRET,required"`
		TokenTTLSeconds int    `env:"JWT_TTL_SECONDS" envDefault:"86400"`
	}

	Admin struct {
		Email    string `env:"ADMIN_EMAIL"`
		Username string `env:"ADMIN_USERNAME"`
		Password string `env:"ADMIN_PASSWORD"`
	}
)

func (p Postgres) URL() string {
	connection := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.User, p.Password),
		Host:   net.JoinHostPort(p.Host, p.Port),
		Path:   p.Database,
	}
	query := connection.Query()
	query.Set("sslmode", p.SSLMode)
	connection.RawQuery = query.Encode()
	return connection.String()
}

func (a Auth) TokenTTL() time.Duration {
	return time.Duration(a.TokenTTLSeconds) * time.Second
}
