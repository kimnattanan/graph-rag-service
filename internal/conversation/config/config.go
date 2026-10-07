package config

import (
	"net"
	"net/url"

	"github.com/kimnattanan/graph-rag-service/internal/common/config"
)

type (
	Config struct {
		Common    config.CommonConfig
		App       App
		Postgres  Postgres
		Knowledge Knowledge
		LLM       LLM
	}

	App struct {
		ServerToRun string `env:"SERVER_TO_RUN,required"`
	}

	Postgres struct {
		Host     string `env:"CONVERSATION_POSTGRES_HOST" envDefault:"localhost"`
		Port     string `env:"CONVERSATION_POSTGRES_PORT" envDefault:"5433"`
		User     string `env:"CONVERSATION_POSTGRES_USER" envDefault:"postgres"`
		Password string `env:"CONVERSATION_POSTGRES_PASSWORD" envDefault:"postgres"`
		Database string `env:"CONVERSATION_POSTGRES_DB" envDefault:"conversation"`
		SSLMode  string `env:"CONVERSATION_POSTGRES_SSLMODE" envDefault:"disable"`
	}

	Knowledge struct {
		GRPCAddress string `env:"KNOWLEDGE_GRPC_ADDRESS" envDefault:"localhost:3000"`
	}

	LLM struct {
		BaseURL string `env:"LLM_BASE_URL" envDefault:""`
		APIKey  string `env:"LLM_API_KEY" envDefault:""`
		Model   string `env:"LLM_MODEL" envDefault:""`
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
