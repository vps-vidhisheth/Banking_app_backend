package service

import (
	"context"
	"errors"

	ledgerService "banking-app/component/ledger/service"
	transactionService "banking-app/component/transactions/service"
	"banking-app/db"
	"banking-app/model"
	"banking-app/repository"

	"github.com/google/uuid"
)

type AccountService struct {
	repo               *repository.Repository[model.Account]
	ledgerService      *ledgerService.LedgerService
	transactionService *transactionService.TransactionService
	db                 *repository.UnitOfWork
}

func NewAccountService(db *repository.UnitOfWork, ledgerSvc *ledgerService.LedgerService, transactionSvc *transactionService.TransactionService) *AccountService {
	return &AccountService{
		repo:               repository.NewRepository[model.Account](db.Tx()),
		ledgerService:      ledgerSvc,
		transactionService: transactionSvc,
		db:                 db,
	}
}

type AccountResponse struct {
	AccountID  uuid.UUID `json:"account_id"`
	BankID     uuid.UUID `json:"bank_id"`
	BankName   string    `json:"bank_name"`
	CustomerID uuid.UUID `json:"customer_id"`
	Balance    float64   `json:"balance"`
	IsActive   bool      `json:"is_active"`
}

func ToAccountResponse(acc *model.Account, bankName string) *AccountResponse {
	return &AccountResponse{
		AccountID:  acc.AccountID,
		BankID:     acc.BankID,
		BankName:   bankName,
		CustomerID: acc.CustomerID,
		Balance:    acc.Balance,
		IsActive:   acc.IsActive,
	}
}

type Bank struct {
	BankID   uuid.UUID `gorm:"type:char(36);primaryKey" json:"bank_id"`
	BankName string    `gorm:"not null" json:"bank_name"`
}

func (s *AccountService) CreateAccountWithUOW(uow *repository.UnitOfWork, customerID, bankID uuid.UUID) (*model.Account, error) {
	acc := &model.Account{
		AccountID:  uuid.New(),
		CustomerID: customerID,
		BankID:     bankID,
		Balance:    0,
		IsActive:   true,
	}
	if err := uow.Tx().Create(acc).Error; err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *AccountService) DepositWithUOW(uow *repository.UnitOfWork, accountID, customerID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}

	ctx := context.Background()

	acc, err := s.repo.WithTransaction(uow.Tx()).GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if acc == nil || acc.CustomerID != customerID {
		return errors.New("caller is not the owner of the account")
	}

	acc.Balance += amount
	if err := uow.Tx().Model(&model.Account{}).
		Where("account_id = ?", acc.AccountID).
		Update("balance", acc.Balance).Error; err != nil {
		return err
	}

	if err := s.transactionService.RecordDeposit(ctx, acc.AccountID, amount, uow.Tx()); err != nil {
		return err
	}

	return nil
}

func (s *AccountService) WithdrawWithUOW(uow *repository.UnitOfWork, accountID, customerID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("withdraw amount must be positive")
	}

	ctx := context.Background()

	acc, err := s.repo.WithTransaction(uow.Tx()).GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if acc == nil || acc.CustomerID != customerID {
		return errors.New("caller is not the owner of the account")
	}
	if acc.Balance < amount {
		return errors.New("insufficient balance")
	}

	acc.Balance -= amount
	if err := uow.Tx().Model(&model.Account{}).
		Where("account_id = ?", acc.AccountID).
		Update("balance", acc.Balance).Error; err != nil {
		return err
	}

	if err := s.transactionService.RecordWithdrawal(ctx, acc.AccountID, amount, uow.Tx()); err != nil {
		return err
	}

	return nil
}

func (s *AccountService) TransferWithUOW(
	uow *repository.UnitOfWork,
	fromAccID, toAccID, fromCustomerID, toCustomerID uuid.UUID,
	amount float64,
) error {
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	ctx := context.Background()

	// Fetch source & destination
	fromAcc, err := s.repo.WithTransaction(uow.Tx()).GetByID(ctx, fromAccID)
	if err != nil {
		return err
	}
	toAcc, err := s.repo.WithTransaction(uow.Tx()).GetByID(ctx, toAccID)
	if err != nil {
		return err
	}

	if fromAcc == nil || toAcc == nil {
		return errors.New("accounts not found")
	}
	if fromAcc.CustomerID != fromCustomerID {
		return errors.New("caller is not owner of source account")
	}
	if fromAcc.Balance < amount {
		return errors.New("insufficient balance")
	}

	fromAcc.Balance -= amount
	toAcc.Balance += amount
	if err := uow.Tx().Save(fromAcc).Error; err != nil {
		return err
	}
	if err := uow.Tx().Save(toAcc).Error; err != nil {
		return err
	}

	if err := s.transactionService.RecordTransfer(ctx, fromAcc.AccountID, toAcc.AccountID, amount, uow.Tx()); err != nil {
		return err
	}

	if fromAcc.BankID != toAcc.BankID {
		if err := s.ledgerService.CreateLedger(ctx, fromAcc.AccountID, amount, "debit",
			"Transfer to "+toAcc.AccountID.String(), &fromAcc.BankID, &toAcc.BankID, uow.Tx()); err != nil {
			return err
		}
		if err := s.ledgerService.CreateLedger(ctx, toAcc.AccountID, amount, "credit",
			"Transfer from "+fromAcc.AccountID.String(), &fromAcc.BankID, &toAcc.BankID, uow.Tx()); err != nil {
			return err
		}
	}

	return nil
}

func (s *AccountService) SoftDeleteAccountWithUOW(uow *repository.UnitOfWork, accountID uuid.UUID) error {
	acc, err := s.repo.GetOne(context.Background(), "account_id = ? AND is_active = ?", accountID, true)
	if err != nil {
		return err
	}
	if acc == nil {
		return errors.New("account not found")
	}
	acc.IsActive = false

	if err := uow.Tx().Save(acc).Error; err != nil {
		return err
	}
	return uow.Tx().Delete(acc).Error
}

func (s *AccountService) RepoGetByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	acc, err := s.repo.GetOne(ctx, "account_id = ? AND is_active = ?", id, true)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, errors.New("account not found")
	}
	return acc, nil
}

func (s *AccountService) GetAccountByIDForStaff(ctx context.Context, accountID, staffID uuid.UUID) (*model.Account, error) {
	acc, err := s.repo.GetOne(ctx, "account_id = ? AND is_active = ? AND staff_id = ?", accountID, true, staffID)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, errors.New("account not found or not assigned to you")
	}
	return acc, nil
}

func (s *AccountService) ListAccountsForStaff(ctx context.Context, staffID uuid.UUID, offset, limit int, searchQuery string) ([]map[string]interface{}, int64, error) {
	dbInstance := db.GetDB()
	query := dbInstance.Model(&model.Account{}).Where("customer_id = ? AND is_active = ?", staffID, true)

	// Apply search filter if provided
	if searchQuery != "" {
		query = query.Where("account_id = ?", searchQuery)
	}

	var accounts []model.Account
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, acc := range accounts {
		var bank model.Bank
		if err := dbInstance.Where("bank_id = ?", acc.BankID).First(&bank).Error; err != nil {
			bank.Name = ""
		}
		accountMap := map[string]interface{}{
			"account_id":  acc.AccountID,
			"customer_id": acc.CustomerID,
			"bank_id":     acc.BankID,
			"bank_name":   bank.Name,
			"balance":     acc.Balance,
			"is_active":   acc.IsActive,
		}
		result = append(result, accountMap)
	}

	return result, total, nil
}

func (s *AccountService) ListAccountsForUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*model.Account, int64, error) {
	var accounts []model.Account
	var total int64

	db := s.repo.GetDB()
	query := db.Where("is_active = 1 AND customer_id = ?", userID)

	if err := query.Model(&model.Account{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	accountPtrs := make([]*model.Account, len(accounts))
	for i := range accounts {
		accountPtrs[i] = &accounts[i]
	}

	return accountPtrs, total, nil
}
