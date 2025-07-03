package auth

import (
	"time"

	"github.com/google/uuid"
)

// models/password_reset_token.go
type PasswordResetToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;" validate:"required,uuid4"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index" validate:"required,uuid4"`
	Token     string    `json:"token" gorm:"type:text;not null" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:timestamptz;not null" validate:"required"`
	Used      bool      `json:"used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null" validate:"required"`
}
