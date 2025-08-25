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
	db.AutoMigrate(&model.Bank{})
	return &BankRepository{db: db}
}

func (r *BankRepository) Create(bank *model.Bank) error {
	if bank == nil {
		return errors.New("bank cannot be nil")
	}
	return r.db.Create(bank).Error
}

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

func (r *BankRepository) List() ([]*model.Bank, error) {
	var banks []*model.Bank
	if err := r.db.Preload("Accounts").Find(&banks).Error; err != nil {
		return nil, err
	}
	return banks, nil
}

func (r *BankRepository) Update(bank *model.Bank) error {
	if bank == nil {
		return errors.New("bank cannot be nil")
	}
	return r.db.Save(bank).Error
}

func (r *BankRepository) Delete(bankID uuid.UUID) error {
	return r.db.Delete(&model.Bank{}, "bank_id = ?", bankID).Error
}
