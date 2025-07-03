package budget

import (
	"time"

	"github.com/google/uuid"
)

// BudgetType represents the frequency of the budget
type BudgetType string

const (
	Weekly  BudgetType = "week"
	Monthly BudgetType = "month"
	Yearly  BudgetType = "year"
)

// models/budget.go
type Budget struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;" validate:"required,uuid4"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index" validate:"required,uuid4"`
	Type      BudgetType `json:"type" gorm:"type:varchar(10);not null" validate:"required,oneof=week month year"`
	Name      string     `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
	StartDate time.Time  `json:"start_date" gorm:"type:timestamp;not null" validate:"required"`
	EndDate   time.Time  `json:"end_date" gorm:"type:timestamp;not null" validate:"required,gtfield=StartDate"`
	Amount    float64    `json:"amount" gorm:"type:decimal(12,2);not null" validate:"required,gt=0"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
