package config

import (
	"github.com/kimnattanan/graph-rag-service/internal/common/config"
)

type (
	Config struct {
		Common   config.CommonConfig
		App      App
		Memgraph Memgraph
		Worker   Worker
	}

	App struct {
		ServerToRun string `env:"SERVER_TO_RUN,required"`
	}

	Memgraph struct {
		Host     string `env:"MEMGRAPH_HOST" envDefault:"localhost"`
		Port     string `env:"MEMGRAPH_PORT" envDefault:"7687"`
		User     string `env:"MEMGRAPH_USER" envDefault:"memgraph"`
		Password string `env:"MEMGRAPH_PASSWORD" envDefault:"memgraph"`
	}

	Worker struct {
		WorkerCount    int `env:"KNOWLEDGE_WORKER_COUNT" envDefault:"0"`
		WorkerInterval int `env:"KNOWLEDGE_WORKER_INTERVAL" envDefault:"60"`
	}
)
