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
func Init() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", config.Config.DBHost, config.Config.DBPort, config.Config.DBUser, config.Config.DBName, config.Config.DBSSLMode)

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
