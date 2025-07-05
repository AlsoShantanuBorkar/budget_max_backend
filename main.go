package main

import (
	"log"

	"github.com/AlsoShantanuBorkar/budget_max/config"
	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/redis"
	"github.com/AlsoShantanuBorkar/budget_max/routes"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	utils.InitializeValidator()
	database.Init()
	redis.InitRedis()
	r := gin.Default()
	api := r.Group("/api/v1")
	routes.RegisterAuthRoutes(api)
	routes.RegisterTransactionRoutes(api)
	routes.RegisterCategoryRoutes(api)
	routes.RegisterBudgetRoutes(api)
	routes.RegisterReportsRoutes(api)
	r.Run(":8080")
}
