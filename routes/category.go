package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(rg *gin.RouterGroup, ctrl controllers.CategoryControllerInterface, sessionDatabaseService database.SessionDatabaseServiceInterface) {
	category := rg.Group("/category")

	// All category routes require authentication
	category.Use(middleware.AuthMiddleware(sessionDatabaseService))

	// GET /api/v1/category - Get all categories for the user
	category.GET("/", ctrl.GetAllCategories)

	// GET /api/v1/category/:id - Get a specific category by ID
	category.GET("/:id", ctrl.GetCategoryByID)

	// POST /api/v1/category - Create a new category
	category.POST("/", ctrl.CreateCategory)

	// PUT /api/v1/category/:id - Update an existing category
	category.PUT("/:id", ctrl.UpdateCategory)

	// DELETE /api/v1/category/:id - Delete a category
	category.DELETE("/:id", ctrl.DeleteCategory)
}
