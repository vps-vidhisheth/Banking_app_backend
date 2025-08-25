package service

import (
	"banking-app/model"
	"banking-app/repository"
	"errors"

	"github.com/google/uuid"
)

type AccountService struct {
	repo               *repository.AccountRepository
	ledgerService      *LedgerService
	transactionService *TransactionService
}

func NewAccountService(
	repo *repository.AccountRepository,
	ledgerService *LedgerService,
	transactionService *TransactionService,
) *AccountService {
	return &AccountService{
		repo:               repo,
		ledgerService:      ledgerService,
		transactionService: transactionService,
	}
}

func (s *AccountService) CreateAccount(customerID, bankID uuid.UUID) (*model.Account, error) {
	acc := &model.Account{
		AccountID:  uuid.New(),
		CustomerID: customerID,
		BankID:     bankID,
		Balance:    0,
	}

	if err := s.repo.Create(acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *AccountService) GetAccountByID(id uuid.UUID) (*model.Account, error) {
	return s.repo.GetByID(id)
}

func (s *AccountService) ListAccounts(page, limit int, customerID, bankID uuid.UUID) ([]*model.Account, error) {
	return s.repo.List(page, limit, customerID, bankID)
}

func (s *AccountService) UpdateAccount(acc *model.Account) error {
	return s.repo.Update(acc)
}

func (s *AccountService) DeleteAccount(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *AccountService) Deposit(accountID, customerID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	acc, err := s.repo.GetByID(accountID)
	if err != nil {
		return err
	}
	if acc.CustomerID != customerID {
		return errors.New("caller is not the owner of the account")
	}

	acc.Balance += amount
	if err := s.repo.Update(acc); err != nil {
		return err
	}

	if err := s.transactionService.RecordDeposit(acc.AccountID, amount); err != nil {
		return err
	}

	_, err = s.ledgerService.CreateLedger(acc.AccountID, amount, "credit", "Deposit", nil, nil)
	return err
}

func (s *AccountService) Withdraw(accountID, customerID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("withdraw amount must be positive")
	}
	acc, err := s.repo.GetByID(accountID)
	if err != nil {
		return err
	}
	if acc.CustomerID != customerID {
		return errors.New("caller is not the owner of the account")
	}
	if acc.Balance < amount {
		return errors.New("insufficient balance")
	}

	acc.Balance -= amount
	if err := s.repo.Update(acc); err != nil {
		return err
	}

	if err := s.transactionService.RecordWithdrawal(acc.AccountID, amount); err != nil {
		return err
	}

	_, err = s.ledgerService.CreateLedger(acc.AccountID, amount, "debit", "Withdraw", nil, nil)
	return err
}

func (s *AccountService) Transfer(fromAccID, toAccID, fromCustomerID, toCustomerID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	fromAcc, err := s.repo.GetByID(fromAccID)
	if err != nil {
		return err
	}
	if fromAcc.CustomerID != fromCustomerID {
		return errors.New("caller is not the owner of the source account")
	}
	if fromAcc.Balance < amount {
		return errors.New("insufficient balance")
	}

	toAcc, err := s.repo.GetByID(toAccID)
	if err != nil {
		return err
	}
	if toAcc.CustomerID != toCustomerID {
		return errors.New("destination account does not belong to the specified customer")
	}

	fromAcc.Balance -= amount
	toAcc.Balance += amount
	if err := s.repo.Update(fromAcc); err != nil {
		return err
	}
	if err := s.repo.Update(toAcc); err != nil {
		return err
	}

	if err := s.transactionService.RecordTransfer(fromAcc.AccountID, toAcc.AccountID, amount); err != nil {
		return err
	}

	_, err = s.ledgerService.CreateLedger(fromAcc.AccountID, amount, "debit", "Transfer to account "+toAccID.String(), &fromAcc.BankID, &toAcc.BankID)
	if err != nil {
		return err
	}
	_, err = s.ledgerService.CreateLedger(toAcc.AccountID, amount, "credit", "Transfer from account "+fromAccID.String(), &fromAcc.BankID, &toAcc.BankID)
	return err
}

func (s *AccountService) CountAccounts(customerID, bankID uuid.UUID) (int64, error) {
	return s.repo.Count(customerID, bankID)
}
