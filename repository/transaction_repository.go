// package repository

// import (
// 	"banking-app/model"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type TransactionRepository struct {
// 	db *gorm.DB
// }

// func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
// 	return &TransactionRepository{db: db}
// }

// func (r *TransactionRepository) Create(transaction *model.Transaction) error {
// 	return r.db.Create(transaction).Error
// }

// func (r *TransactionRepository) GetByAccountID(accountID uuid.UUID) ([]model.Transaction, error) {
// 	var transactions []model.Transaction
// 	err := r.db.Where("account_id = ?", accountID).Find(&transactions).Error
// 	return transactions, err
// }

// func (r *TransactionRepository) GetNetTransfer(accountID uuid.UUID) (float64, error) {
// 	var sent float64
// 	var received float64

// 	err := r.db.Model(&model.Transaction{}).
// 		Where("account_id = ? AND type = ?", accountID, model.Transfer).
// 		Select("COALESCE(SUM(amount),0)").Scan(&sent).Error
// 	if err != nil {
// 		return 0, err
// 	}

// 	err = r.db.Model(&model.Transaction{}).
// 		Where("related_account_id = ? AND type = ?", accountID, model.Transfer).
// 		Select("COALESCE(SUM(amount),0)").Scan(&received).Error
// 	if err != nil {
// 		return 0, err
// 	}

// 	return received - sent, nil
// }

// func (r *TransactionRepository) GetAll() ([]model.Transaction, error) {
// 	var transactions []model.Transaction
// 	err := r.db.Find(&transactions).Error
// 	return transactions, err
// }

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

// ✅ NEW: Paginated GetAll
func (r *TransactionRepository) GetAllPaginated(limit, offset int) ([]model.Transaction, int64, error) {
	var transactions []model.Transaction
	var total int64

	// count total rows
	if err := r.db.Model(&model.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// fetch paginated
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// ✅ NEW: Paginated GetByAccountID
func (r *TransactionRepository) GetByAccountIDPaginated(accountID uuid.UUID, limit, offset int) ([]model.Transaction, int64, error) {
	var transactions []model.Transaction
	var total int64

	// count total for this account
	if err := r.db.Model(&model.Transaction{}).
		Where("account_id = ?", accountID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// fetch paginated
	err := r.db.Where("account_id = ?", accountID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}
