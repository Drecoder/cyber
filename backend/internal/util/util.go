package util

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitLogger initializes the global Logger and returns a cleanup function
func InitLogger() func() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// Return cleanup function to flush logs at program exit
	return func() {
		Logger.Sync()
	}
}
