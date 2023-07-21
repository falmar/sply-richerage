package zaplogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

func New(debug bool) *zap.Logger {
	var lConfig zap.Config

	if os.Getenv("DEBUG") != "" || debug {
		lConfig = zap.NewDevelopmentConfig()
		lConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		lConfig = zap.NewProductionConfig()
	}

	lConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if strings.ToUpper(os.Getenv("LOG_LEVEL")) == "DEBUG" {
		lConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, _ := lConfig.Build()

	return logger
}
