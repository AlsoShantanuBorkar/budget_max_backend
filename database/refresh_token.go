package database

import (
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
)

func CreateRefreshToken(token *models.RefreshToken) error {
	return DB.Create(token).Error
}

func GetRefreshTokenByToken(token uuid.UUID) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := DB.First(&refreshToken, "token = ?", token).Error
	return &refreshToken, err
}

func RevokeRefreshToken(tokenID uuid.UUID) error {
	return DB.Model(&models.RefreshToken{}).Where("id = ?", tokenID).Update("revoked", true).Error
}
func RevokeRefreshTokensBySessionID(sessionID uuid.UUID) error {
	return DB.Model(&models.RefreshToken{}).Where("session_id = ?", sessionID).Update("revoked", true).Error
}

func DeleteRefreshToken(tokenID uuid.UUID) error {
	return DB.Delete(&models.RefreshToken{}, "id = ?", tokenID).Error
}
