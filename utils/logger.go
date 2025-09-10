package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	AppLogger         *zap.Logger
	AuthLogger        *zap.Logger
	BudgetLogger      *zap.Logger
	CategoryLogger    *zap.Logger
	ReportLogger      *zap.Logger
	TransactionLogger *zap.Logger
	TwoFactorLogger   *zap.Logger	
	DatabaseLogger    *zap.Logger
)

func newLogger(file string) *zap.Logger {
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/" + file,
		MaxSize:    10,   // MB
		MaxBackups: 5,    // keep last 5
		MaxAge:     30,   // days
		Compress:   true, // gzip old logs
	})

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)),
		zap.DebugLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}
func WithRequestID(logger *zap.Logger, requestID string) *zap.Logger {
	return logger.With(zap.String("request_id", requestID))
}

func InitLogger() {
	AppLogger = newLogger("app.log")
	AuthLogger = newLogger("auth.log")
	TransactionLogger = newLogger("transactions.log")
	ReportLogger = newLogger("reports.log")
	CategoryLogger = newLogger("categories.log")
	BudgetLogger = newLogger("budgets.log")
	TwoFactorLogger = newLogger("2fa.log")
	DatabaseLogger = newLogger("database.log")
	AppLogger.Info("Logger initialized")
}

