package config

import "github.com/kimnattanan/graph-rag-service/internal/common/config"

type (
	Config struct {
		Common config.CommonConfig
		App App
	}

	App struct {
		ServerToRun string `env:"SERVER_TO_RUN,required"`
	}
)