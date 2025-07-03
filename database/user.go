package database

import (
	"errors"
	"strings"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}

func GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := DB.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := DB.Debug().First(&user, "LOWER(email) = LOWER(?)", strings.TrimSpace(email)).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func UpdateUser(userID uuid.UUID, updates map[string]interface{}) error {
	return DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func DeleteUser(id uuid.UUID) error {
	return DB.Delete(&models.User{}, "id = ?", id).Error
}
