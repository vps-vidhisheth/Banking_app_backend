package service

import (
	"banking-app/model"
	"banking-app/repository"
	"errors"

	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}

// --- Record deposit, withdrawal, transfer ---
func (s *TransactionService) RecordDeposit(accountID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	tx := &model.Transaction{
		TransactionID: uuid.New(),
		AccountID:     accountID,
		Amount:        amount,
		Type:          model.Deposit,
		Note:          "Deposit",
	}
	return s.transactionRepo.Create(tx)
}

func (s *TransactionService) RecordWithdrawal(accountID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	tx := &model.Transaction{
		TransactionID: uuid.New(),
		AccountID:     accountID,
		Amount:        amount,
		Type:          model.Withdraw,
		Note:          "Withdrawal",
	}
	return s.transactionRepo.Create(tx)
}

func (s *TransactionService) RecordTransfer(fromID, toID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	txFrom := &model.Transaction{
		TransactionID:    uuid.New(),
		AccountID:        fromID,
		RelatedAccountID: &toID,
		Amount:           amount,
		Type:             model.Transfer,
		Note:             "Transfer to account",
	}
	if err := s.transactionRepo.Create(txFrom); err != nil {
		return err
	}

	txTo := &model.Transaction{
		TransactionID:    uuid.New(),
		AccountID:        toID,
		RelatedAccountID: &fromID,
		Amount:           amount,
		Type:             model.Transfer,
		Note:             "Transfer from account",
	}
	return s.transactionRepo.Create(txTo)
}

// --- Fetch transactions for a specific account ---
func (s *TransactionService) GetTransactionsByAccount(accountID uuid.UUID) ([]model.Transaction, error) {
	return s.transactionRepo.GetByAccountID(accountID)
}

// --- Fetch net transfers ---
func (s *TransactionService) GetNetTransfers(accountID uuid.UUID) (float64, error) {
	return s.transactionRepo.GetNetTransfer(accountID)
}

// --- Fetch all transactions (for admin, optional filtering in future) ---
func (s *TransactionService) GetAllTransactions() ([]model.Transaction, error) {
	// You can implement pagination/filtering in repo later
	return s.transactionRepo.GetAll()
}
