package database

import (
	"fmt"
	"log"

	"github.com/AlsoShantanuBorkar/budget_max/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// update with your actual module path
)

var database *gorm.DB

var budgetDBService *BudgetService
var categoryDBService *CategoryDatabaseService
var transactionDBService *TransactionDatabaseService
var userDBService *UserDatabaseService
var sessionDBService *SessionDatabaseService
var refreshTokenDBService RefreshTokenDatabaseService

// Init loads env, connects to DB, runs migrations
func Init() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", config.Config.DBHost, config.Config.DBPort, config.Config.DBUser, config.Config.DBName, config.Config.DBSSLMode)

	// Open DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}
	database = db

	var result string
	database.Raw("SELECT current_database();").Scan(&result)

	log.Printf("Connected to database: %s", result)

	// Initialize services
	budgetDBService = &BudgetService{database: database}
	categoryDBService = &CategoryDatabaseService{database: database}
	transactionDBService = &TransactionDatabaseService{database: database}
	userDBService = &UserDatabaseService{database: database}
	sessionDBService = &SessionDatabaseService{database: database}
	refreshTokenDBService = RefreshTokenDatabaseService{database: database}

	log.Println("Database services initialized")
}
