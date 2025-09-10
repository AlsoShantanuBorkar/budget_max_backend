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

func NewConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	Config := &AppConfig{
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

	return Config, nil
}
