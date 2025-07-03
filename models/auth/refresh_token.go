package auth

import (
	"time"

	"github.com/google/uuid"
)

// models/refresh_token.go
type RefreshToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;" validate:"required,uuid4"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index" validate:"required,uuid4"`
	SessionID uuid.UUID `json:"session_id" gorm:"type:uuid;not null;index" validate:"required,uuid4"`
	Token     uuid.UUID `json:"token" gorm:"type:uuid;uniqueIndex;not null" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:timestamptz;not null" validate:"required"`
	Revoked   bool      `json:"revoked" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null" validate:"required"`
}
