package auth

import (
	"time"

	"github.com/google/uuid"
)

// models/session.go
type Session struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;" validate:"required,uuid4"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index" validate:"required,uuid4"`
	Token     uuid.UUID `json:"token" gorm:"type:uuid;uniqueIndex;not null" validate:"required"`
	UserAgent string    `json:"user_agent" gorm:"type:text"`
	IPAddress string    `json:"ip_address" gorm:"type:varchar(45)"`
	Revoked   bool      `gorm:"default:false"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:timestamptz;not null" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null" validate:"required"`
}
