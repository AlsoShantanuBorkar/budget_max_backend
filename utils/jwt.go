package utils

import (
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/config"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/models/auth"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user *models.User) (string, error) {
	claims := &auth.TwoFAClaims{
		UserID: user.ID,
		Email:  user.Email,
		Is2FA:  true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "budget_max",
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.JWTSecret))
}

func VerifyJWT(token string) (bool, *models.TwoFAClaims, error) {
	claims := &models.TwoFAClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.JWTSecret), nil
	})

	if err != nil {
		return false, nil, err
	}

	return parsedToken.Valid, claims, nil
}
