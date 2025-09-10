package database

import (
	"errors"
	"strings"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(userID uuid.UUID, updates map[string]interface{}) error
	DeleteUser(id uuid.UUID) error
}	
type UserDatabaseService struct {
	database *gorm.DB
}
func NewUserDatabaseService(db *gorm.DB) UserRepository {
	return &UserDatabaseService{database: db}
}

func (s *UserDatabaseService) CreateUser(user *models.User) error {
	return s.database.Create(user).Error
}

func (s *UserDatabaseService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.database.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (s *UserDatabaseService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := s.database.Debug().First(&user, "LOWER(email) = LOWER(?)", strings.TrimSpace(email)).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (s *UserDatabaseService) UpdateUser(userID uuid.UUID, updates map[string]interface{}) error {
	return s.database.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (s *UserDatabaseService) DeleteUser(id uuid.UUID) error {
	return s.database.Delete(&models.User{}, "id = ?", id).Error
}
