package repository

import (
	"banking-app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

// Create a new account
func (r *AccountRepository) Create(acc *model.Account) error {
	// Generate UUID if not already set
	if acc.AccountID == uuid.Nil {
		acc.AccountID = uuid.New()
	}
	return r.DB.Create(acc).Error
}

// Get account by UUID
func (r *AccountRepository) GetByID(id uuid.UUID) (*model.Account, error) {
	var acc model.Account
	if err := r.DB.First(&acc, "account_id = ? AND is_active = ?", id, true).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}

// Update account
func (r *AccountRepository) Update(acc *model.Account) error {
	return r.DB.Save(acc).Error
}

// Soft delete by UUID
func (r *AccountRepository) Delete(id uuid.UUID) error {
	return r.DB.Model(&model.Account{}).
		Where("account_id = ?", id).
		Update("is_active", false).Error
}

// List accounts with optional filtering
func (r *AccountRepository) List(limit, offset int, customerID, bankID uuid.UUID) ([]*model.Account, error) {
	var accounts []*model.Account
	query := r.DB.Model(&model.Account{}).Where("is_active = ?", true)

	if customerID != uuid.Nil {
		query = query.Where("customer_id = ?", customerID)
	}
	if bankID != uuid.Nil {
		query = query.Where("bank_id = ?", bankID)
	}

	if err := query.Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) Count(customerID, bankID uuid.UUID) (int64, error) {
	var count int64
	query := r.DB.Model(&model.Account{}).Where("is_active = ?", true)
	if customerID != uuid.Nil {
		query = query.Where("customer_id = ?", customerID)
	}
	if bankID != uuid.Nil {
		query = query.Where("bank_id = ?", bankID)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
