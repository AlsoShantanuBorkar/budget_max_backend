package services

import (
	"net/http"
	"sort"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models/reports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetBudgetSummary(c *gin.Context, budgetID uuid.UUID, userId uuid.UUID) (*reports.BudgetSummary, *ServiceError) {
	budget, err := database.GetBudgetByID(budgetID, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error finding budget")
	}
	txns, err := database.GetTransactionsByBudget(userId, budgetID)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error finding Transactions")
	}
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.BudgetSummary{
		BudgetID:      budgetID,
		BudgetName:    budget.Name,
		BudgetAmount:  budget.Amount,
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func GetWeeklySummary(c *gin.Context, userId uuid.UUID, startDate time.Time, endDate time.Time) (*reports.WeeklySummary, *ServiceError) {
	txns, err := database.GetTransactionsByDateRange(userId, startDate, endDate)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error fetching transactions")
	}
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.WeeklySummary{
		StartDate:     startDate.Format("2006-01-02"),
		EndDate:       endDate.Format("2006-01-02"),
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func GetMonthlySummary(c *gin.Context, userId uuid.UUID, month time.Time) (*reports.MonthlySummary, *ServiceError) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	txns, err := database.GetTransactionsByDateRange(userId, startOfMonth, endOfMonth)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error fetching transactions")
	}
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.MonthlySummary{
		Month:         month.Format("January 2006"),
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func GetYearlySummary(c *gin.Context, userId uuid.UUID, year time.Time) (*reports.YearlySummary, *ServiceError) {
	startOfYear := time.Date(year.Year(), 1, 1, 0, 0, 0, 0, year.Location())
	endOfYear := time.Date(year.Year(), 12, 31, 23, 59, 59, 999999999, year.Location())

	txns, err := database.GetTransactionsByDateRange(userId, startOfYear, endOfYear)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error fetching transactions")
	}
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.YearlySummary{
		Year:          year.Format("2006"),
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func GetCategorySummary(c *gin.Context, userId uuid.UUID, categoryID uuid.UUID) (*reports.CategorySummary, *ServiceError) {
	category, err := database.GetCategoryByID(categoryID, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error fetching category")
	}
	if category == nil {
		return nil, NewServiceError(http.StatusNotFound, "Category not found")
	}

	txns, err := database.GetTransactionsByCategory(userId, categoryID)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error fetching transactions")
	}
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.CategorySummary{
		CategoryID:    categoryID,
		CategoryName:  category.Name,
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func GetCustomDateRangeSummary(c *gin.Context, userId uuid.UUID, startDate time.Time, endDate time.Time) (*reports.CustomDateRangeSummary, *ServiceError) {
	txns, err := database.GetTransactionsByDateRange(userId, startDate, endDate)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error fetching transactions")
	}
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.CustomDateRangeSummary{
		StartDate:     startDate.Format("2006-01-02"),
		EndDate:       endDate.Format("2006-01-02"),
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func GetDailyAverageSummary(c *gin.Context, userId uuid.UUID, startDate time.Time, endDate time.Time) (*reports.DailyAverageSummary, *ServiceError) {
	txns, err := database.GetTransactionsByDateRange(userId, startDate, endDate)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error fetching transactions")
	}

	totalExpenses := 0.0
	totalIncome := 0.0
	daysDiff := int(endDate.Sub(startDate).Hours()/24) + 1

	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}

	dailyExpenses := totalExpenses / float64(daysDiff)
	dailyIncome := totalIncome / float64(daysDiff)
	dailyNetBalance := dailyIncome - dailyExpenses

	return &reports.DailyAverageSummary{
		Date:          startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02"),
		TotalExpenses: dailyExpenses,
		TotalIncome:   dailyIncome,
		NetBalance:    dailyNetBalance,
	}, nil
}

func GetTopCategories(c *gin.Context, userId uuid.UUID, limit int, transactionType string) ([]*reports.TopCategory, *ServiceError) {
	if limit <= 0 {
		limit = 5
	}

	// Get all transactions for the user
	txns, err := database.GetTransactionsByType(userId, transactionType)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error fetching transactions")
	}

	// Group transactions by category
	categoryTotals := make(map[string]float64)
	categoryNames := make(map[string]string)

	for _, txn := range txns {
		categoryTotals[txn.CategoryID.String()] += txn.Amount
		// We'll get category names in the next step
	}

	// Get category names
	for categoryID := range categoryTotals {
		parsedId, err := uuid.Parse(categoryID)
		if err != nil {
			continue
		}
		category, err := database.GetCategoryByID(parsedId, userId)
		if err != nil {
			continue
		}
		if category != nil {
			categoryNames[categoryID] = category.Name
		}
	}

	// Convert to slice for sorting
	var categories []*reports.TopCategory
	for categoryID, amount := range categoryTotals {
		categoryName, exists := categoryNames[categoryID]
		if !exists {
			categoryName = "Unknown Category"
		}
		categories = append(categories, &reports.TopCategory{
			CategoryID:   categoryID,
			CategoryName: categoryName,
			Amount:       amount,
			Type:         transactionType,
		})
	}

	// Sort by amount (descending for expenses, ascending for income)
	if transactionType == "expense" {
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].Amount > categories[j].Amount
		})
	} else {
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].Amount > categories[j].Amount
		})
	}

	// Limit results and add rank
	if len(categories) > limit {
		categories = categories[:limit]
	}

	for i, category := range categories {
		category.Rank = i + 1
	}

	return categories, nil
}

func GetAllCategoriesSummary(c *gin.Context, userId uuid.UUID) ([]*reports.CategorySummary, *ServiceError) {
	categories, err := database.GetUserCategories(userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "Error fetching categories")
	}

	var summaries []*reports.CategorySummary

	for _, category := range categories {
		txns, err := database.GetTransactionsByCategory(userId, category.ID)
		if err != nil {
			continue
		}

		totalExpenses := 0.0
		totalIncome := 0.0
		for _, txn := range txns {
			if txn.Type == "expense" {
				totalExpenses += txn.Amount
			} else {
				totalIncome += txn.Amount
			}
		}

		netBalance := totalIncome - totalExpenses
		summaries = append(summaries, &reports.CategorySummary{
			CategoryID:    category.ID,
			CategoryName:  category.Name,
			TotalExpenses: totalExpenses,
			TotalIncome:   totalIncome,
			NetBalance:    netBalance,
		})
	}

	return summaries, nil
}
