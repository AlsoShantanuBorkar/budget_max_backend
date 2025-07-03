package categories

import "github.com/google/uuid"

// models/category.go
type Category struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;" validate:"required,uuid4"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index" validate:"required,uuid4"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null" validate:"required"`
	Type      string    `json:"type" gorm:"type:varchar(10);not null" validate:"required,oneof=expense income"`
	Icon      *string   `json:"icon" gorm:"type:varchar(50);default:null"` // pointer string for nullable
	IsDefault bool      `json:"is_default" gorm:"default:false"`
}
