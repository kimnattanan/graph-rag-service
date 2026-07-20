package logs

import (
	"github.com/kimnattanan/graph-rag-service/internal/common/config"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func Init(cfg *config.CommonConfig) {
	SetFormatter(logrus.StandardLogger(), cfg)
	logrus.SetLevel(logrus.DebugLevel)
}

func SetFormatter(logger *logrus.Logger, cfg *config.CommonConfig) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})

	if cfg.IsLocalEnv {
		logger.SetFormatter(&prefixed.TextFormatter{
			ForceFormatting: true,
		})
	}
}
