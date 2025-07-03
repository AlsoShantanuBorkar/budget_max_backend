package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(rg *gin.RouterGroup) {
	category := rg.Group("/category")

	// All category routes require authentication
	category.Use(middleware.AuthMiddleware())

	// GET /api/v1/category - Get all categories for the user
	category.GET("/", controllers.GetAllCategories)

	// GET /api/v1/category/:id - Get a specific category by ID
	category.GET("/:id", controllers.GetCategoryByID)

	// POST /api/v1/category - Create a new category
	category.POST("/", controllers.CreateCategory)

	// PUT /api/v1/category/:id - Update an existing category
	category.PUT("/:id", controllers.UpdateCategory)

	// DELETE /api/v1/category/:id - Delete a category
	category.DELETE("/:id", controllers.DeleteCategory)
}
