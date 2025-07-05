package reports

import "github.com/google/uuid"

type BudgetSummary struct {
	BudgetID      uuid.UUID `json:"budget_id"`
	BudgetName    string    `json:"budget_name"`
	BudgetAmount  float64   `json:"budget_amount"`
	TotalExpenses float64   `json:"total_expenses"`
	TotalIncome   float64   `json:"total_income"`
	NetBalance    float64   `json:"net_balance"`
}

type WeeklySummary struct {
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	TotalExpenses float64 `json:"total_expenses"`
	TotalIncome   float64 `json:"total_income"`
	NetBalance    float64 `json:"net_balance"`
}

type MonthlySummary struct {
	Month         string  `json:"month"`
	TotalExpenses float64 `json:"total_expenses"`
	TotalIncome   float64 `json:"total_income"`
	NetBalance    float64 `json:"net_balance"`
}

type YearlySummary struct {
	Year          string  `json:"year"`
	TotalExpenses float64 `json:"total_expenses"`
	TotalIncome   float64 `json:"total_income"`
	NetBalance    float64 `json:"net_balance"`
}

type CategorySummary struct {
	CategoryID    uuid.UUID `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	TotalExpenses float64   `json:"total_expenses"`
	TotalIncome   float64   `json:"total_income"`
	NetBalance    float64   `json:"net_balance"`
}

type CustomDateRangeSummary struct {
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	TotalExpenses float64 `json:"total_expenses"`
	TotalIncome   float64 `json:"total_income"`
	NetBalance    float64 `json:"net_balance"`
}
type DailyAverageSummary struct {
	Date          string  `json:"date"`
	TotalExpenses float64 `json:"total_expenses"`
	TotalIncome   float64 `json:"total_income"`
	NetBalance    float64 `json:"net_balance"`
}

type TopCategory struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Type         string  `json:"type"` // income or expense
	Rank         int     `json:"rank"`
}
