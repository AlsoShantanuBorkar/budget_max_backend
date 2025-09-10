package main

import (
	"log"

	"github.com/AlsoShantanuBorkar/budget_max/config"
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/redis"
	"github.com/AlsoShantanuBorkar/budget_max/routes"
	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize configuration, database, Redis, and validator
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	utils.InitializeValidator()
	db, err := database.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	redis.InitRedis()
	r := gin.Default()
	api := r.Group("/api/v1")

	// Initialize Database Services
	userDatabaseService := database.NewUserDatabaseService(db)
	sessionDatabaseService := database.NewSessionDatabaseService(db)
	categoryDatabaseService := database.NewCategoryDatabaseService(db)
	budgetDatabaseService := database.NewBudgetDatabaseService(db)
	transactionDatabaseService := database.NewTransactionDatabaseService(db)
	refreshTokenDatabaseService := database.NewRefreshTokenDatabaseService(db)

	// Initialize Services
	authService := services.NewAuthService(userDatabaseService, sessionDatabaseService, refreshTokenDatabaseService)
	budgetService := services.NewBudgetService(budgetDatabaseService)
	categoryService := services.NewCategoryService(categoryDatabaseService)
	transactionService := services.NewTransactionService(transactionDatabaseService)
	reportsService := services.NewReportsService(transactionDatabaseService, categoryDatabaseService, budgetDatabaseService)

	// Initialize Controllers
	authController := controllers.NewAuthController(authService)
	budgetController := controllers.NewBudgetController(budgetService)
	categoryController := controllers.NewCategoryController(categoryService)
	transactionController := controllers.NewTransactionController(transactionService)
	reportsController := controllers.NewReportsController(reportsService)

	// Register Routes

	// Unprotected Routes
	routes.RegisterUnprotectedAuthRoutes(api, authController)

	// Protected Routes
	routes.RegisterAuthRoutes(api, authController, sessionDatabaseService)
	routes.RegisterTransactionRoutes(api, transactionController, sessionDatabaseService)
	routes.RegisterCategoryRoutes(api, categoryController, sessionDatabaseService)
	routes.RegisterBudgetRoutes(api, budgetController, sessionDatabaseService)
	routes.RegisterReportsRoutes(api, reportsController, sessionDatabaseService)
	r.Run(":8080")
}
