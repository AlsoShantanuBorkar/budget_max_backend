package database

import (
	"fmt"
	"log"

	"github.com/AlsoShantanuBorkar/budget_max/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// update with your actual module path
)

// Init loads env, connects to DB, runs migrations
func Init(config *config.AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", config.DBHost, config.DBPort, config.DBUser, config.DBName, config.DBSSLMode)

	// Open DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err

	}

	var result string
	db.Raw("SELECT current_database();").Scan(&result)

	log.Printf("Connected to database: %s", result)

	return db, nil
}
