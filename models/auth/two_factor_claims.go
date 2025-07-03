package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TwoFAClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Is2FA  bool      `json:"is_2fa"`
	jwt.RegisteredClaims
}
