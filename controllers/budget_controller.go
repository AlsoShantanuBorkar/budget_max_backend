package controllers

import (
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateBudget(c *gin.Context) {
	var req models.CreateBudgetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	budget, serviceErr := services.CreateBudget(c, &req, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Budget created successfully",
		"data": gin.H{
			"budget": budget,
		},
	})
}

func UpdateBudget(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	var req models.UpdateBudgetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	budgetIdStr := c.Param("id")
	budgetId, err := uuid.Parse(budgetIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid budget ID format",
		})
		return
	}

	updatedBudget, serviceErr := services.UpdateBudget(c, &req, budgetId, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Budget updated successfully",
		"data":    updatedBudget,
	})
}

func DeleteBudget(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	budgetIdStr := c.Param("id")
	budgetId, err := uuid.Parse(budgetIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid budget ID format",
		})
		return
	}

	serviceErr := services.DeleteBudget(c, budgetId, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Budget deleted successfully",
	})
}

func GetBudgetsByUserID(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	budgets, serviceErr := services.GetBudgetsByUserID(c, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Budgets fetched successfully",
		"data":    budgets,
	})
}

func GetBudgetByID(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	budgetIDStr := c.Param("id")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid budget ID format",
		})
		return
	}

	budget, serviceErr := services.GetBudgetByID(c, budgetID, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Budget fetched successfully",
		"data":    budget,
	})
}
