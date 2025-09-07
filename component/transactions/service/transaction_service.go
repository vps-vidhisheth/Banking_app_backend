package service

import (
	"banking-app/model"
	"banking-app/repository"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	repo *repository.Repository[model.Transaction]
	db   *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		repo: repository.NewRepository[model.Transaction](db),
		db:   db,
	}
}

func (s *TransactionService) GetTransactions(ctx context.Context, accountID *uuid.UUID, txType, note string, startDate, endDate *time.Time, limit, offset int) ([]model.Transaction, int64, error) {
	filters := map[string]interface{}{}

	if accountID != nil && *accountID != uuid.Nil {
		filters["account_id = ?"] = *accountID
	}
	if txType != "" {
		filters["type = ?"] = strings.ToLower(txType)
	}
	if note != "" {
		filters["note LIKE ?"] = "%" + strings.ToLower(note) + "%"
	}

	dbQuery := s.db.WithContext(ctx).Model(&model.Transaction{})
	for k, v := range filters {
		dbQuery = dbQuery.Where(k, v)
	}
	if startDate != nil {
		dbQuery = dbQuery.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		dbQuery = dbQuery.Where("created_at <= ?", *endDate)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, errors.New("no transactions found")
	}

	var results []model.Transaction
	if err := dbQuery.Limit(limit).Offset(offset).Order("created_at DESC").Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, id uuid.UUID) error {
	tx, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if tx == nil {
		return errors.New("transaction not found")
	}
	return nil
}

func (s *TransactionService) RecordDeposit(ctx context.Context, accountID uuid.UUID, amount float64, tx *gorm.DB) error {
	t := &model.Transaction{
		TransactionID: uuid.New(),
		AccountID:     accountID,
		Type:          "deposit",
		Amount:        amount,
		CreatedAt:     time.Now(),
	}
	if tx != nil {
		return tx.Create(t).Error
	}
	return s.db.WithContext(ctx).Create(t).Error
}

func (s *TransactionService) RecordWithdrawal(ctx context.Context, accountID uuid.UUID, amount float64, tx *gorm.DB) error {
	t := &model.Transaction{
		TransactionID: uuid.New(),
		AccountID:     accountID,
		Type:          "withdraw",
		Amount:        amount,
		CreatedAt:     time.Now(),
	}
	if tx != nil {
		return tx.Create(t).Error
	}
	return s.db.WithContext(ctx).Create(t).Error
}

func (s *TransactionService) RecordTransfer(ctx context.Context, fromAccID, toAccID uuid.UUID, amount float64, tx *gorm.DB) error {
	t := &model.Transaction{
		TransactionID: uuid.New(),
		AccountID:     fromAccID,
		Type:          "transfer",
		Amount:        amount,
		CreatedAt:     time.Now(),
	}
	if tx != nil {
		if err := tx.Create(t).Error; err != nil {
			return err
		}
		t2 := &model.Transaction{
			TransactionID: uuid.New(),
			AccountID:     toAccID,
			Type:          "transfer",
			Amount:        amount,
			CreatedAt:     time.Now(),
		}
		return tx.Create(t2).Error
	}

	if err := s.db.WithContext(ctx).Create(t).Error; err != nil {
		return err
	}
	t2 := &model.Transaction{
		TransactionID: uuid.New(),
		AccountID:     toAccID,
		Type:          "transfer",
		Amount:        amount,
		CreatedAt:     time.Now(),
	}
	return s.db.WithContext(ctx).Create(t2).Error
}
