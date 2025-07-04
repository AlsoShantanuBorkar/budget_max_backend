package database

import (
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
)

func CreateTransation(txn *models.Transaction) error {
	return DB.Create(txn).Error
}

func GetTransactionsByUser(userID uuid.UUID) ([]*models.Transaction, error) {
	var txns []*models.Transaction
	err := DB.Where("user_id = ?", userID).Find(&txns).Error
	return txns, err
}

func UpdateTransaction(id uuid.UUID, updates map[string]any) error {
	return DB.Model(&models.Transaction{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteTransaction(id uuid.UUID) error {
	return DB.Delete(&models.Transaction{}, "id = ?", id).Error
}

func GetTransactionByID(txnId uuid.UUID, userId uuid.UUID) (*models.Transaction, error) {
	var txn models.Transaction
	err := DB.First(&txn, "id = ? AND user_id = ?", txnId, userId).Error
	if err != nil {
		return nil, err
	}
	return &txn, err
}
