package database

import (
	"fmt"
	"log"

	"github.com/AlsoShantanuBorkar/budget_max/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// update with your actual module path
)

var DB *gorm.DB

// Init loads env, connects to DB, runs migrations
func Init() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", config.Config.DBHost, config.Config.DBPort, config.Config.DBUser, config.Config.DBName, config.Config.DBSSLMode)

	// Open DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}
	DB = db

	var result string
	DB.Raw("SELECT current_database();").Scan(&result)
}

// func migrate() {
// 	// Add all your models here
// 	err := DB.AutoMigrate()
// 	if err != nil {
// 		log.Fatalf("Failed to migrate database: %v", err)
// 	}
// }
