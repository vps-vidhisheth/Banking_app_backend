package repository

import (
	"banking-app/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LedgerRepository struct {
	db *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) *LedgerRepository {
	return &LedgerRepository{db: db}
}

func (r *LedgerRepository) Create(ledger *model.Ledger) error {
	return r.db.Create(ledger).Error
}

func (r *LedgerRepository) GetByID(id uuid.UUID) (*model.Ledger, error) {
	var ledger model.Ledger
	if err := r.db.First(&ledger, "ledger_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &ledger, nil
}

func (r *LedgerRepository) ListByAccount(accountID uuid.UUID, limit, offset int) ([]model.Ledger, int64, error) {
	var ledgers []model.Ledger
	var total int64

	query := r.db.Model(&model.Ledger{}).Where("account_id = ?", accountID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&ledgers).Error; err != nil {
		return nil, 0, err
	}

	return ledgers, total, nil
}

func (r *LedgerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Ledger{}, "ledger_id = ?", id).Error
}

func (r *LedgerRepository) NetBankTransfer(bankFromID, bankToID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.Model(&model.Ledger{}).
		Select("SUM(amount)").
		Where("bank_from_id = ? AND bank_to_id = ?", bankFromID, bankToID).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
