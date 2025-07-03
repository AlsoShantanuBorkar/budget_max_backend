package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RedisHost     string
	RedisPassword string
	JWTSecret     string
	DBSSLMode     string
}

var Config *AppConfig

func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	Config = &AppConfig{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     os.Getenv("DB_SSLMODE"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
	}

	return nil
}
