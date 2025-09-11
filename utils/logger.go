package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)




var loggerInstance *zerolog.Logger

func InitLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logDirectory := fmt.Sprintf("logs/%s", time.Now().Format("2006-01-02"))
	if err:= os.MkdirAll(logDirectory, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v\n", err))		
	}
	logPath := filepath.Join(logDirectory, "app.log")

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // megabytes
		MaxBackups: 0,  // number of backups
		MaxAge:     30, // days
		Compress:   false,
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	
	multi := zerolog.MultiLevelWriter(consoleWriter, lumberjackLogger)

	l := zerolog.New(multi).With().Timestamp().Logger()
	loggerInstance = &l

	return loggerInstance
}

func GetLogger() *zerolog.Logger {
	if loggerInstance == nil {
		 return InitLogger()
	}
	return loggerInstance
}

func StartDailyRotation() {
    go func() {
        currentDay := time.Now().Day()
        for {
            time.Sleep(time.Minute)
            if time.Now().Day() != currentDay {
                InitLogger()
                currentDay = time.Now().Day()
            }
        }
    }()
}