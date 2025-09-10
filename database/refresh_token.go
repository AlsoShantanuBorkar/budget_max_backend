package database

import (
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {	
	CreateRefreshToken(token *models.RefreshToken) error
	GetRefreshTokenByToken(token uuid.UUID) (*models.RefreshToken, error)
	RevokeRefreshToken(tokenID uuid.UUID) error
	RevokeRefreshTokensBySessionID(sessionID uuid.UUID) error
	DeleteRefreshToken(tokenID uuid.UUID) error
}

type RefreshTokenDatabaseService struct {
	database *gorm.DB
}

func NewRefreshTokenDatabaseService(db *gorm.DB) RefreshTokenRepository {
	return &RefreshTokenDatabaseService{database: db}
}

func (s *RefreshTokenDatabaseService) CreateRefreshToken(token *models.RefreshToken) error {
	return s.database.Create(token).Error
}

func (s *RefreshTokenDatabaseService) GetRefreshTokenByToken(token uuid.UUID) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := s.database.First(&refreshToken, "token = ?", token).Error
	return &refreshToken, err
}

func (s *RefreshTokenDatabaseService) RevokeRefreshToken(tokenID uuid.UUID) error {
	return s.database.Model(&models.RefreshToken{}).Where("id = ?", tokenID).Update("revoked", true).Error
}
func (s *RefreshTokenDatabaseService) RevokeRefreshTokensBySessionID(sessionID uuid.UUID) error {
	return s.database.Model(&models.RefreshToken{}).Where("session_id = ?", sessionID).Update("revoked", true).Error
}

func (s *RefreshTokenDatabaseService) DeleteRefreshToken(tokenID uuid.UUID) error {
	return s.database.Delete(&models.RefreshToken{}, "id = ?", tokenID).Error
}
