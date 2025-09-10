package database

import (
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionDatabaseServiceInterface interface {
	CreateSession(session *models.Session) error
	GetSessionByToken(token uuid.UUID) (*models.Session, error)
	DeleteSession(sessionID uuid.UUID) error
	RevokeSession(tokenID uuid.UUID) error
}

type SessionDatabaseService struct {
	database *gorm.DB
}

func NewSessionDatabaseService(db *gorm.DB) SessionDatabaseServiceInterface {
	return &SessionDatabaseService{database: db}
}

func (s *SessionDatabaseService) CreateSession(session *models.Session) error {
	return s.database.Create(session).Error
}

func (s *SessionDatabaseService) GetSessionByToken(token uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := s.database.Where("token = ? AND revoked = false", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionDatabaseService) DeleteSession(sessionID uuid.UUID) error {
	return s.database.Delete(&models.Session{}, "id = ?", sessionID).Error
}

func (s *SessionDatabaseService) RevokeSession(tokenID uuid.UUID) error {
	return s.database.Model(&models.Session{}).Where("id = ?", tokenID).Update("revoked", true).Error
}
