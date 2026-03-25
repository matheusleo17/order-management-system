package logger

import (
	"os"

	"go.uber.org/zap"
)

func New() *zap.Logger {
	if os.Getenv("APP_ENV") == "development" {
		log, err := zap.NewDevelopment()
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}
		return log
	}

	log, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	return log
}
