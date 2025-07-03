package database

import (
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
)

func CreateSession(session *models.Session) error {
	return DB.Create(session).Error
}

func GetSessionByToken(token uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := DB.Where("token = ? AND revoked = false", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func DeleteSession(sessionID uuid.UUID) error {
	return DB.Delete(&models.Session{}, "id = ?", sessionID).Error
}

func RevokeSession(tokenID uuid.UUID) error {
	return DB.Model(&models.Session{}).Where("id = ?", tokenID).Update("revoked", true).Error
}
