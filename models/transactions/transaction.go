package transactions

import (
	"time"

	"github.com/google/uuid"
)

// models/transaction.go
type Transaction struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey" validate:"required,uuid4"`
	UserID     uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index" validate:"required,uuid4"`
	Amount     float64    `json:"amount" gorm:"not null" validate:"required,gt=0"`
	Type       string     `json:"type" gorm:"type:varchar(10);not null" validate:"required,oneof=expense income"`
	Name       string     `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
	Date       time.Time  `json:"date" gorm:"type:timestamptz;not null" validate:"required"`
	Note       string     `json:"note" gorm:"type:text"`
	CategoryID *uuid.UUID `json:"category_id,omitempty" gorm:"type:uuid" validate:"omitempty,uuid4"`
	CreatedAt  time.Time  `json:"created_at" gorm:"type:timestamptz;not null" validate:"required"`
}
