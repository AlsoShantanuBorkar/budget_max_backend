package utils

import (
	"errors"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateSessionAndRefreshToken(user *models.User, ip, agent string, sessionDatabaseService database.SessionDatabaseServiceInterface, refreshTokenDatabaseService database.RefreshTokenDatabaseServiceInterface, c *gin.Context) (models.Session, models.RefreshToken, error) {

	sessionToken := uuid.New()
	refreshToken := uuid.New()
	now := time.Now()

	session := models.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: now,
		ExpiresAt: now.Add(15 * time.Minute),
		IPAddress: ip,
		UserAgent: agent,
		Token:     sessionToken,
	}

	refresh := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		SessionID: session.ID,
		Token:     refreshToken,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
		Revoked:   false,
	}

	if err := sessionDatabaseService.CreateSession(&session); err != nil {
		return models.Session{}, models.RefreshToken{}, errors.New("failed to create session")
	}

	if err := refreshTokenDatabaseService.CreateRefreshToken(&refresh); err != nil {
		return models.Session{}, models.RefreshToken{}, errors.New("failed to create refresh token")
	}

	return session, refresh, nil
}
