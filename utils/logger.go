package utils

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger() {
	var err error

	Logger, err = zap.NewDevelopment()
	if err != nil {
		panic("failed to init logger: " + err.Error())

	}
	defer Logger.Sync()
}
