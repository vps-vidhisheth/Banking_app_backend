package service

import (
	"banking-app/model"
	"banking-app/repository"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LedgerService struct {
	repo *repository.LedgerRepository
}

func NewLedgerService(repo *repository.LedgerRepository) *LedgerService {
	return &LedgerService{repo: repo}
}

func (s *LedgerService) CreateLedger(accountID uuid.UUID, amount float64, ledgerType, description string, bankFromID, bankToID *uuid.UUID) (*model.Ledger, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	entryType := strings.ToLower(strings.TrimSpace(ledgerType))
	if entryType != "debit" && entryType != "credit" {
		return nil, errors.New("ledgerType must be 'debit' or 'credit'")
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

	if err := s.repo.Create(entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *LedgerService) GetAllLedgers(accountID uuid.UUID, limit, offset int) ([]model.Ledger, int64, error) {
	return s.repo.ListByAccount(accountID, limit, offset)
}

func (s *LedgerService) DeleteLedger(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *LedgerService) GetNetBankTransfer(bankFromID, bankToID uuid.UUID) (float64, error) {
	return s.repo.NetBankTransfer(bankFromID, bankToID)
}

func (s *LedgerService) GetLedger(id uuid.UUID) (*model.Ledger, error) {
	return s.repo.GetByID(id)
}
