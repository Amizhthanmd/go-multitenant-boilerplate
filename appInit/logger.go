package appInit

import "go.uber.org/zap"

func ZapLogger(Mode string) *zap.Logger {
	logger := zap.Must(zap.NewProduction())
	if Mode == "debug" {
		logger = zap.Must(zap.NewDevelopment())
	}
	return logger
}
