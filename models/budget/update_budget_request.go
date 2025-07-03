package budget

import "time"

type UpdateBudgetRequest struct {
	Name      *string    `json:"name,omitempty" validate:"omitempty"`
	StartDate *time.Time `json:"start_date,omitempty" validate:"omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty" validate:"omitempty,gtfield=StartDate"`
	Amount    *float64   `json:"amount,omitempty" validate:"omitempty,gt=0"`
}
