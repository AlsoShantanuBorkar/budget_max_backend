package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateBudget(c *gin.Context) {
	var req models.CreateBudgetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	if err := utils.GetValidator().Struct(req); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	b := models.Budget{
		ID:        uuid.New(),
		UserID:    userId,
		Type:      req.Type,
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}

	err := database.CreateBudget(&b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create budget",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Budget created successfully",
		"data": gin.H{
			"budget": b,
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

	// Fetch existing budget
	_, err = database.GetBudgetByID(budgetId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "budget not found",
		})
		return
	}

	// Update only fields that are set
	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = *req.EndDate
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}

	// Save updated budget
	err = database.UpdateBudget(budgetId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update budget",
		})
		return
	}

	// Fetch updated budget
	updatedBudget, err := database.GetBudgetByID(budgetId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch updated budget",
		})
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

	_, err = database.GetBudgetByID(budgetId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "budget not found",
		})
		return
	}

	err = database.DeleteBudget(budgetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while deleting budget",
		})
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

	budgets, err := database.GetBudgetsByUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while fetching budgets",
		})
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

	b, err := database.GetBudgetByID(budgetID, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "budget not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Budget fetched successfully",
		"data":    b,
	})
}
