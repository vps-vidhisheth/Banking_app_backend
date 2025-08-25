package repository

import (
	"banking-app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) GetByAccountID(accountID uuid.UUID) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("account_id = ?", accountID).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetNetTransfer(accountID uuid.UUID) (float64, error) {
	var sent float64
	var received float64

	err := r.db.Model(&model.Transaction{}).
		Where("account_id = ? AND type = ?", accountID, model.Transfer).
		Select("COALESCE(SUM(amount),0)").Scan(&sent).Error
	if err != nil {
		return 0, err
	}

	err = r.db.Model(&model.Transaction{}).
		Where("related_account_id = ? AND type = ?", accountID, model.Transfer).
		Select("COALESCE(SUM(amount),0)").Scan(&received).Error
	if err != nil {
		return 0, err
	}

	return received - sent, nil
}

func (r *TransactionRepository) GetAll() ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Find(&transactions).Error
	return transactions, err
}
