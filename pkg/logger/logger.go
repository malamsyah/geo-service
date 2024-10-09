package logger

import (
	"go.uber.org/zap"
)

// nolint: gochecknoglobals
var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func Errorf(format string, v ...interface{}) {
	logger.Sugar().Errorf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Sugar().Debugf(format, v...)
}
