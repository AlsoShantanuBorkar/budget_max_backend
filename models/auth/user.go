package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id" validate:"required,uuid4"`
	Email            string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password         string    `gorm:"not null" json:"-" validate:"required,min=8"`
	TwoFactorEnabled bool      `gorm:"default:false" json:"two_factor_enabled"`
	TwoFactorSecret  string    `gorm:"" json:"-"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at" validate:"required,datetime"`
}
