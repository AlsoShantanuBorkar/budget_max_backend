package utils

import (
	"os"

	"github.com/rs/zerolog"
)




var loggerInstance *zerolog.Logger

func InitLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	l := zerolog.New(file).With().Timestamp().Logger()
	loggerInstance = &l
	return loggerInstance
}

func GetLogger() *zerolog.Logger {
	if loggerInstance == nil {
		 return InitLogger()
	}
	return loggerInstance
}
