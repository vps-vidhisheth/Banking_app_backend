package repository

import (
	"banking-app/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankRepository struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) *BankRepository {
	db.AutoMigrate(&model.Bank{}) // optional: migrate the Bank table
	return &BankRepository{db: db}
}

// Create a new bank
func (r *BankRepository) Create(bank *model.Bank) error {
	if bank == nil {
		return errors.New("bank cannot be nil")
	}
	return r.db.Create(bank).Error
}

// GetByID fetches a bank by UUID
func (r *BankRepository) GetByID(bankID uuid.UUID) (*model.Bank, error) {
	var bank model.Bank
	if err := r.db.Preload("Accounts").First(&bank, "bank_id = ?", bankID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bank, nil
}

// List returns all banks
func (r *BankRepository) List() ([]*model.Bank, error) {
	var banks []*model.Bank
	if err := r.db.Preload("Accounts").Find(&banks).Error; err != nil {
		return nil, err
	}
	return banks, nil
}

// Update updates an existing bank
func (r *BankRepository) Update(bank *model.Bank) error {
	if bank == nil {
		return errors.New("bank cannot be nil")
	}
	return r.db.Save(bank).Error
}

// Delete deletes a bank by UUID
func (r *BankRepository) Delete(bankID uuid.UUID) error {
	return r.db.Delete(&model.Bank{}, "bank_id = ?", bankID).Error
}
