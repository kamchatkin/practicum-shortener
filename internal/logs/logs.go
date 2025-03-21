package logs

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func NewLogger() *zap.Logger {
	return logger
}
