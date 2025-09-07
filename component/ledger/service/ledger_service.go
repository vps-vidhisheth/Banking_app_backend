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

type LedgerService struct {
	repo *repository.Repository[model.Ledger]
	db   *gorm.DB
}

func NewLedgerService(db *gorm.DB) *LedgerService {
	return &LedgerService{
		repo: repository.NewRepository[model.Ledger](db),
		db:   db,
	}
}

func (s *LedgerService) CreateLedger(ctx context.Context, accountID uuid.UUID, amount float64, ledgerType, description string, bankFromID, bankToID *uuid.UUID, tx *gorm.DB) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	entryType := strings.ToLower(strings.TrimSpace(ledgerType))
	if entryType != "debit" && entryType != "credit" {
		return errors.New("ledgerType must be 'debit' or 'credit'")
	}

	var txType string
	if bankFromID != nil && bankToID != nil {
		txType = "transfer"
	} else if entryType == "credit" {
		txType = "deposit"
	} else {
		txType = "withdraw"
	}

	entry := &model.Ledger{
		LedgerID:        uuid.New(),
		AccountID:       &accountID,
		BankFromID:      bankFromID,
		BankToID:        bankToID,
		Amount:          amount,
		TransactionType: txType,
		EntryType:       entryType,
		Description:     description,
		CreatedAt:       time.Now(),
	}

	if tx != nil {
		return s.repo.WithTransaction(tx).Create(ctx, entry)
	}

	return s.db.WithContext(ctx).Transaction(func(t *gorm.DB) error {
		return s.repo.WithTransaction(t).Create(ctx, entry)
	})
}

func (s *LedgerService) DeleteLedger(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.repo.WithTransaction(tx).DeleteByID(ctx, id)
	})
}

func (s *LedgerService) CheckLedgerExists(ctx context.Context, id uuid.UUID) error {
	ledger, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ledger == nil {
		return errors.New("ledger not found")
	}
	return nil
}

func (s *LedgerService) CheckLedgersWithFilters(ctx context.Context, filters map[string]string) error {
	query := make(map[string]interface{})

	for key, value := range filters {
		if value == "" {
			continue
		}
		switch key {
		case "account_id":
			id, err := uuid.Parse(value)
			if err != nil {
				return errors.New("invalid account_id")
			}
			query["account_id = ?"] = id
		case "entry_type":
			ltype := strings.ToLower(strings.TrimSpace(value))
			if ltype == "debit" || ltype == "credit" {
				query["entry_type = ?"] = ltype
			}
		case "transaction_type":
			ttype := strings.ToLower(strings.TrimSpace(value))
			if ttype == "deposit" || ttype == "withdraw" || ttype == "transfer" {
				query["transaction_type = ?"] = ttype
			}
		}
	}

	ledgers, err := s.repo.List(ctx, 1, 0, query)
	if err != nil {
		return err
	}
	if len(ledgers) == 0 {
		return errors.New("no ledgers found with given filters")
	}
	return nil
}

func (s *LedgerService) CheckAnyLedgers(ctx context.Context) error {
	ledgers, err := s.repo.List(ctx, 1, 0, nil)
	if err != nil {
		return err
	}
	if len(ledgers) == 0 {
		return errors.New("no ledgers found")
	}
	return nil
}

func (s *LedgerService) GetAllLedgers(ctx context.Context, accountID *uuid.UUID, limit, offset int) ([]model.Ledger, int64, error) {
	filters := map[string]interface{}{}
	if accountID != nil && *accountID != uuid.Nil {
		filters["account_id = ?"] = *accountID
	}

	ledgers, err := s.repo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}
	if len(ledgers) == 0 {
		return nil, 0, errors.New("no ledgers found")
	}

	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return ledgers, count, nil
}

func (s *LedgerService) GetNetBankTransfer(ctx context.Context, bankFromID, bankToID uuid.UUID) (float64, error) {
	var debitSum, reverseDebitSum float64

	if err := s.db.WithContext(ctx).
		Model(&model.Ledger{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("bank_from_id = ? AND bank_to_id = ? AND entry_type = ?", bankFromID, bankToID, "debit").
		Scan(&debitSum).Error; err != nil {
		return 0, err
	}

	if err := s.db.WithContext(ctx).
		Model(&model.Ledger{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("bank_from_id = ? AND bank_to_id = ? AND entry_type = ?", bankToID, bankFromID, "debit").
		Scan(&reverseDebitSum).Error; err != nil {
		return 0, err
	}

	netTransfer := debitSum - reverseDebitSum

	return netTransfer, nil
}
