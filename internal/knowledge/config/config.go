package config

import "github.com/kimnattanan/graph-rag-service/internal/common/config"

type (
	Config struct {
		Common   config.CommonConfig
		App      App
		Memgraph Memgraph
	}

	App struct {
		ServerToRun string `env:"SERVER_TO_RUN,required"`
	}

	Memgraph struct {
		Host     string `env:"MEMGRAPH_HOST,default=localhost"`
		Port     string `env:"MEMGRAPH_PORT,default=7687"`
		User     string `env:"MEMGRAPH_USER,default=memgraph"`
		Password string `env:"MEMGRAPH_PASSWORD,default=memgraph"`
	}
)
