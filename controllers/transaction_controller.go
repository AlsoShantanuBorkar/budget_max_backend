package controllers

import (
	"net/http"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTransation(c *gin.Context) {

	var req models.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	userIdRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User ID not found",
		})
		return
	}

	userId, ok := userIdRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user ID type",
		})
		return
	}

	parsedDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Date",
		})
		return
	}

	txn := models.Transaction{
		ID:         uuid.New(),
		UserID:     userId,
		Amount:     req.Amount,
		Type:       req.Type,
		Note:       req.Note,
		Name:       req.Name,
		Date:       parsedDate,
		CategoryID: req.CategoryIDs,
		CreatedAt:  time.Now(),
	}

	err = database.CreateTransation(&txn)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create transaction",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction created successfully",
		"data": gin.H{
			"transaction": txn,
		},
	})

}

func UpdateTransaction(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	var req models.UpdateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	if err := utils.GetValidator().Struct(req); err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	txnIdStr := c.Param("id")
	txnId, err := uuid.Parse(txnIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid transaction ID format",
		})
		return
	}

	// Fetch existing transaction
	_, err = database.GetTransactionByID(txnId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "transaction not found",
		})
		return
	}

	updates := make(map[string]any)
	if req.Name != nil {

		updates["name"] = *req.Name
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Date != nil {

		parsedDate, err := time.Parse(time.RFC3339, *req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Date",
			})
			return
		}
		updates["date"] = parsedDate
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.CategoryID != nil {
		updates["category_id"] = req.CategoryID
	}

	// Save updated transaction
	err = database.UpdateTransaction(txnId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update transaction",
		})
		return
	}

	updatedTransaction, err := database.GetTransactionByID(txnId, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch updated budget",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction updated successfully",
		"data":    updatedTransaction,
	})

}

func DeleteTransaction(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	txnIdStr := c.Param("id")
	txnId, err := uuid.Parse(txnIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid transaction ID format",
		})
		return
	}

	_, err = database.GetTransactionByID(txnId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "transaction not found",
		})
		return
	}

	err = database.DeleteTransaction(txnId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Occurred",
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Transaction deleted successfully",
	})

}

func GetTransactionsByUserID(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	txns, err := database.GetTransactionsByUser(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Occurred while fetching transactions",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions fetched successfully",
		"data":    txns,
	})

}

func GetTransactionByID(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	txnIDStr := c.Param("id")
	txnID, err := uuid.Parse(txnIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid transaction ID format",
		})
		return
	}

	txn, err := database.GetTransactionByID(txnID, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction fetched successfully",
		"data":    txn,
	})

}
