package budget

import "time"

type CreateBudgetRequest struct {
	Type      BudgetType `json:"type" validate:"required,oneof=week month year"`
	Name      string     `json:"name" validate:"required"`
	StartDate time.Time  `json:"start_date" validate:"required"`
	EndDate   time.Time  `json:"end_date" validate:"required,gtfield=StartDate"`
	Amount    float64    `json:"amount" validate:"required,gt=0"`
}
