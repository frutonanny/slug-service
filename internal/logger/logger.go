package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func Must() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("zap new development: %v", err))
	}
	defer func() {
		_ = logger.Sync()
	}()

	return logger
}
