// package service

// import (
// 	ledgerService "banking-app/component/ledger/service"
// 	transactionService "banking-app/component/transactions/service"
// 	"banking-app/model"
// 	"banking-app/repository"
// 	"banking-app/utils"
// 	"context"
// 	"errors"
// 	"net/http"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type AccountService struct {
// 	repo               *repository.Repository[model.Account]
// 	ledgerService      *ledgerService.LedgerService
// 	transactionService *transactionService.TransactionService
// 	db                 *gorm.DB
// }

// func NewAccountService(
// 	db *gorm.DB,
// 	ledgerSvc *ledgerService.LedgerService,
// 	transactionSvc *transactionService.TransactionService,
// ) *AccountService {
// 	return &AccountService{
// 		repo:               repository.NewRepository[model.Account](db),
// 		ledgerService:      ledgerSvc,
// 		transactionService: transactionSvc,
// 		db:                 db,
// 	}
// }

// // ---------------- Create / Get ----------------
// func (s *AccountService) CreateAccount(ctx context.Context, customerID, bankID uuid.UUID) (*model.Account, error) {
// 	acc := &model.Account{
// 		AccountID:  uuid.New(),
// 		CustomerID: customerID,
// 		BankID:     bankID,
// 		Balance:    0,
// 		IsActive:   true,
// 	}

// 	if err := s.repo.Create(ctx, acc); err != nil {
// 		return nil, err
// 	}
// 	return acc, nil
// }

// func (s *AccountService) GetAccountByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
// 	return s.repo.GetOne(ctx, "account_id = ? AND is_active = ?", id, true)
// }

// // ---------------- Pagination ----------------
// func (s *AccountService) ListAccountsWithPagination(ctx context.Context, r *http.Request, customerID, bankID uuid.UUID) (map[string]interface{}, error) {
// 	pagination := utils.GetPaginationParams(r, 10, 0)
// 	offset := pagination.Offset
// 	limit := pagination.Limit

// 	filters := map[string]interface{}{"is_active = ?": true}
// 	if customerID != uuid.Nil {
// 		filters["customer_id = ?"] = customerID
// 	}
// 	if bankID != uuid.Nil {
// 		filters["bank_id = ?"] = bankID
// 	}

// 	accounts, err := s.repo.List(ctx, limit, offset, filters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := s.repo.Count(ctx, filters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return utils.PaginatedResponse(accounts, total, limit, offset), nil
// }

// // ---------------- Deposit ----------------
// func (s *AccountService) Deposit(ctx context.Context, accountID, customerID uuid.UUID, amount float64) error {
// 	if amount <= 0 {
// 		return errors.New("deposit amount must be positive")
// 	}

// 	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		acc, err := s.GetAccountByID(ctx, accountID)
// 		if err != nil {
// 			return err
// 		}
// 		if acc == nil || acc.CustomerID != customerID {
// 			return errors.New("caller is not the owner of the account")
// 		}

// 		acc.Balance += amount
// 		if err := tx.Save(acc).Error; err != nil {
// 			return err
// 		}

// 		if err := s.transactionService.RecordDeposit(ctx, acc.AccountID, amount, tx); err != nil {
// 			return err
// 		}

// 		if _, err := s.ledgerService.CreateLedger(ctx, acc.AccountID, amount, "credit", "Deposit", nil, nil, tx); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

// // ---------------- Withdraw ----------------
// func (s *AccountService) Withdraw(ctx context.Context, accountID, customerID uuid.UUID, amount float64) error {
// 	if amount <= 0 {
// 		return errors.New("withdraw amount must be positive")
// 	}

// 	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		acc, err := s.GetAccountByID(ctx, accountID)
// 		if err != nil {
// 			return err
// 		}
// 		if acc == nil || acc.CustomerID != customerID {
// 			return errors.New("caller is not the owner of the account")
// 		}
// 		if acc.Balance < amount {
// 			return errors.New("insufficient balance")
// 		}

// 		acc.Balance -= amount
// 		if err := tx.Save(acc).Error; err != nil {
// 			return err
// 		}

// 		if err := s.transactionService.RecordWithdrawal(ctx, acc.AccountID, amount, tx); err != nil {
// 			return err
// 		}

// 		if _, err := s.ledgerService.CreateLedger(ctx, acc.AccountID, amount, "debit", "Withdraw", nil, nil, tx); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

// // ---------------- Transfer ----------------
// func (s *AccountService) Transfer(ctx context.Context, fromAccID, toAccID, fromCustomerID, toCustomerID uuid.UUID, amount float64) error {
// 	if amount <= 0 {
// 		return errors.New("transfer amount must be positive")
// 	}

// 	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		fromAcc, err := s.GetAccountByID(ctx, fromAccID)
// 		if err != nil {
// 			return err
// 		}
// 		if fromAcc == nil || fromAcc.CustomerID != fromCustomerID {
// 			return errors.New("caller is not the owner of the source account")
// 		}
// 		if fromAcc.Balance < amount {
// 			return errors.New("insufficient balance")
// 		}

// 		toAcc, err := s.GetAccountByID(ctx, toAccID)
// 		if err != nil {
// 			return err
// 		}
// 		if toAcc == nil || toAcc.CustomerID != toCustomerID {
// 			return errors.New("destination account does not belong to the specified customer")
// 		}

// 		fromAcc.Balance -= amount
// 		toAcc.Balance += amount

// 		if err := tx.Save(fromAcc).Error; err != nil {
// 			return err
// 		}
// 		if err := tx.Save(toAcc).Error; err != nil {
// 			return err
// 		}

// 		if err := s.transactionService.RecordTransfer(ctx, fromAcc.AccountID, toAcc.AccountID, amount, tx); err != nil {
// 			return err
// 		}

// 		if _, err := s.ledgerService.CreateLedger(ctx, fromAcc.AccountID, amount, "debit", "Transfer to account "+toAccID.String(), &fromAcc.BankID, &toAcc.BankID, tx); err != nil {
// 			return err
// 		}
// 		if _, err := s.ledgerService.CreateLedger(ctx, toAcc.AccountID, amount, "credit", "Transfer from account "+fromAccID.String(), &fromAcc.BankID, &toAcc.BankID, tx); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

// // ---------------- Count ----------------
// func (s *AccountService) CountAccounts(ctx context.Context, customerID, bankID uuid.UUID) (int64, error) {
// 	filters := map[string]interface{}{"is_active = ?": true}
// 	if customerID != uuid.Nil {
// 		filters["customer_id = ?"] = customerID
// 	}
// 	if bankID != uuid.Nil {
// 		filters["bank_id = ?"] = bankID
// 	}
// 	return s.repo.Count(ctx, filters)
// }

// // ---------------- Update ----------------
// func (s *AccountService) UpdateAccount(ctx context.Context, acc *model.Account) error {
// 	if acc == nil {
// 		return errors.New("account is nil")
// 	}
// 	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Save(acc).Error; err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }

// // ---------------- Soft Delete ----------------
// func (s *AccountService) SoftDeleteAccount(ctx context.Context, accountID uuid.UUID) error {
// 	acc, err := s.GetAccountByID(ctx, accountID)
// 	if err != nil {
// 		return err
// 	}
// 	if acc == nil {
// 		return errors.New("account not found")
// 	}

// 	// Mark as inactive
// 	acc.IsActive = false

// 	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// Save the IsActive update first
// 		if err := tx.Save(acc).Error; err != nil {
// 			return err
// 		}

// 		// Then soft delete using GORM's Delete to set DeletedAt
// 		if err := tx.Delete(acc).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

package service

import (
	"context"
	"errors"
	"net/http"

	ledgerService "banking-app/component/ledger/service"
	transactionService "banking-app/component/transactions/service"
	"banking-app/model"
	"banking-app/repository"
	"banking-app/utils"

	"github.com/google/uuid"
)

type AccountService struct {
	repo               *repository.Repository[model.Account]
	ledgerService      *ledgerService.LedgerService
	transactionService *transactionService.TransactionService
	db                 *repository.UnitOfWork
}

func NewAccountService(
	db *repository.UnitOfWork,
	ledgerSvc *ledgerService.LedgerService,
	transactionSvc *transactionService.TransactionService,
) *AccountService {
	return &AccountService{
		repo:               repository.NewRepository[model.Account](db.Tx()),
		ledgerService:      ledgerSvc,
		transactionService: transactionSvc,
		db:                 db,
	}
}

// ---------------- Create Account ----------------
func (s *AccountService) CreateAccountWithUOW(uow *repository.UnitOfWork, customerID, bankID uuid.UUID) error {
	acc := &model.Account{
		AccountID:  uuid.New(),
		CustomerID: customerID,
		BankID:     bankID,
		Balance:    0,
		IsActive:   true,
	}
	return uow.Tx().Create(acc).Error
}

// ---------------- Deposit ----------------
func (s *AccountService) DepositWithUOW(uow *repository.UnitOfWork, accountID, customerID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}

	ctx := context.Background()

	// Fetch latest account inside the transaction
	acc, err := s.repo.WithTransaction(uow.Tx()).GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if acc == nil || acc.CustomerID != customerID {
		return errors.New("caller is not the owner of the account")
	}

	// Update balance
	acc.Balance += amount
	if err := uow.Tx().Model(&model.Account{}).
		Where("account_id = ?", acc.AccountID).
		Update("balance", acc.Balance).Error; err != nil {
		return err
	}

	// Record transaction and ledger inside the same tx
	if err := s.transactionService.RecordDeposit(ctx, acc.AccountID, amount, uow.Tx()); err != nil {
		return err
	}
	if err := s.ledgerService.CreateLedger(ctx, acc.AccountID, amount, "credit", "Deposit", nil, nil, uow.Tx()); err != nil {
		return err
	}

	return nil
}

// ---------------- Withdraw ----------------
func (s *AccountService) WithdrawWithUOW(uow *repository.UnitOfWork, accountID, customerID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("withdraw amount must be positive")
	}

	ctx := context.Background()

	// Fetch latest account inside the transaction
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

	// Update balance
	acc.Balance -= amount
	if err := uow.Tx().Model(&model.Account{}).
		Where("account_id = ?", acc.AccountID).
		Update("balance", acc.Balance).Error; err != nil {
		return err
	}

	// Record transaction and ledger inside the same tx
	if err := s.transactionService.RecordWithdrawal(ctx, acc.AccountID, amount, uow.Tx()); err != nil {
		return err
	}
	if err := s.ledgerService.CreateLedger(ctx, acc.AccountID, amount, "debit", "Withdraw", nil, nil, uow.Tx()); err != nil {
		return err
	}

	return nil
}

// ---------------- Transfer ----------------
func (s *AccountService) TransferWithUOW(uow *repository.UnitOfWork, fromAccID, toAccID, fromCustomerID, toCustomerID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	ctx := context.Background()

	// Fetch latest source account
	fromAcc, err := s.repo.WithTransaction(uow.Tx()).GetByID(ctx, fromAccID)
	if err != nil {
		return err
	}
	if fromAcc == nil || fromAcc.CustomerID != fromCustomerID {
		return errors.New("caller is not the owner of the source account")
	}
	if fromAcc.Balance < amount {
		return errors.New("insufficient balance")
	}

	// Fetch latest destination account
	toAcc, err := s.repo.WithTransaction(uow.Tx()).GetByID(ctx, toAccID)
	if err != nil {
		return err
	}
	if toAcc == nil || toAcc.CustomerID != toCustomerID {
		return errors.New("destination account does not belong to the specified customer")
	}

	// Update balances
	fromAcc.Balance -= amount
	toAcc.Balance += amount
	if err := uow.Tx().Model(&model.Account{}).Where("account_id = ?", fromAcc.AccountID).Update("balance", fromAcc.Balance).Error; err != nil {
		return err
	}
	if err := uow.Tx().Model(&model.Account{}).Where("account_id = ?", toAcc.AccountID).Update("balance", toAcc.Balance).Error; err != nil {
		return err
	}

	// Record transfer transaction
	if err := s.transactionService.RecordTransfer(ctx, fromAcc.AccountID, toAcc.AccountID, amount, uow.Tx()); err != nil {
		return err
	}

	// Ledger entries
	if err := s.ledgerService.CreateLedger(ctx, fromAcc.AccountID, amount, "debit", "Transfer to account "+toAccID.String(), &fromAcc.BankID, &toAcc.BankID, uow.Tx()); err != nil {
		return err
	}
	if err := s.ledgerService.CreateLedger(ctx, toAcc.AccountID, amount, "credit", "Transfer from account "+fromAccID.String(), &fromAcc.BankID, &toAcc.BankID, uow.Tx()); err != nil {
		return err
	}

	return nil
}

// ---------------- Update ----------------
func (s *AccountService) UpdateAccountWithUOW(uow *repository.UnitOfWork, acc *model.Account) error {
	if acc == nil {
		return errors.New("account is nil")
	}

	// Only update specific fields, including BankID
	return uow.Tx().Model(&model.Account{}).
		Where("account_id = ?", acc.AccountID).
		Updates(map[string]interface{}{
			"customer_id": acc.CustomerID,
			"bank_id":     acc.BankID,
			"balance":     acc.Balance,
			"is_active":   acc.IsActive,
		}).Error
}

// ---------------- Soft Delete ----------------
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

// ---------------- Repo Get ----------------
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

// ---------------- List Accounts with Pagination ----------------
func (s *AccountService) ListAccountsWithPagination(ctx context.Context, r *http.Request) ([]*model.Account, error) {
	pagination := utils.GetPaginationParams(r, 10, 0)
	offset := pagination.Offset
	limit := pagination.Limit

	filters := map[string]interface{}{"is_active": true}

	query := r.URL.Query()
	if customer := query.Get("customer_id"); customer != "" {
		custUUID, err := uuid.Parse(customer)
		if err == nil {
			filters["customer_id = ?"] = custUUID
		}
	}
	if bank := query.Get("bank_id"); bank != "" {
		bankUUID, err := uuid.Parse(bank)
		if err == nil {
			filters["bank_id = ?"] = bankUUID
		}
	}

	accounts, err := s.repo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, err
	}

	accountPtrs := make([]*model.Account, len(accounts))
	for i := range accounts {
		accountPtrs[i] = &accounts[i]
	}
	return accountPtrs, nil
}

// ---------------- List Accounts by UserID ----------------
func (s *AccountService) ListAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Account, error) {
	accounts, err := s.repo.List(ctx, 0, 0, map[string]interface{}{
		"is_active = ?":   true,
		"customer_id = ?": userID,
	})
	if err != nil {
		return nil, err
	}

	accountPtrs := make([]*model.Account, len(accounts))
	for i := range accounts {
		accountPtrs[i] = &accounts[i]
	}
	return accountPtrs, nil
}

func (s *AccountService) ListAccountsWithPaginationByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*model.Account, int64, error) {
	filters := map[string]interface{}{
		"customer_id = ?": userID,
		"is_active = ?":   true,
	}

	// Get total count
	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated accounts
	accounts, err := s.repo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}

	accountPtrs := make([]*model.Account, len(accounts))
	for i := range accounts {
		accountPtrs[i] = &accounts[i]
	}

	return accountPtrs, total, nil
}
