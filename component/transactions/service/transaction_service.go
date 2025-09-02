// package service

// import (
// 	"banking-app/model"
// 	"banking-app/repository"
// 	"banking-app/utils"
// 	"context"
// 	"errors"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type TransactionService struct {
// 	repo *repository.Repository[model.Transaction]
// 	db   *gorm.DB
// }

// func NewTransactionService(db *gorm.DB) *TransactionService {
// 	return &TransactionService{
// 		repo: repository.NewRepository[model.Transaction](db),
// 		db:   db,
// 	}
// }

// // RecordDeposit supports optional transaction
// func (s *TransactionService) RecordDeposit(ctx context.Context, accountID uuid.UUID, amount float64, tx *gorm.DB) error {
// 	if amount <= 0 {
// 		return errors.New("deposit amount must be positive")
// 	}

// 	entry := &model.Transaction{
// 		TransactionID: uuid.New(),
// 		AccountID:     accountID,
// 		Amount:        amount,
// 		Type:          model.Deposit,
// 		Note:          "Deposit",
// 	}

// 	repo := s.repo
// 	if tx != nil {
// 		repo = repo.WithTransaction(tx)
// 	}
// 	return repo.Create(ctx, entry)
// }

// // RecordWithdrawal supports optional transaction
// func (s *TransactionService) RecordWithdrawal(ctx context.Context, accountID uuid.UUID, amount float64, tx *gorm.DB) error {
// 	if amount <= 0 {
// 		return errors.New("withdrawal amount must be positive")
// 	}

// 	entry := &model.Transaction{
// 		TransactionID: uuid.New(),
// 		AccountID:     accountID,
// 		Amount:        amount,
// 		Type:          model.Withdraw,
// 		Note:          "Withdrawal",
// 	}

// 	repo := s.repo
// 	if tx != nil {
// 		repo = repo.WithTransaction(tx)
// 	}
// 	return repo.Create(ctx, entry)
// }

// // RecordTransfer supports optional transaction for atomic debit/credit
// func (s *TransactionService) RecordTransfer(ctx context.Context, fromID, toID uuid.UUID, amount float64, tx *gorm.DB) error {
// 	if amount <= 0 {
// 		return errors.New("transfer amount must be positive")
// 	}

// 	repo := s.repo
// 	if tx != nil {
// 		repo = repo.WithTransaction(tx)
// 	}

// 	// Debit from source
// 	txFrom := &model.Transaction{
// 		TransactionID:    uuid.New(),
// 		AccountID:        fromID,
// 		RelatedAccountID: &toID,
// 		Amount:           amount,
// 		Type:             model.Transfer,
// 		Note:             "Transfer to account",
// 	}
// 	if err := repo.Create(ctx, txFrom); err != nil {
// 		return err
// 	}

// 	// Credit to destination
// 	txTo := &model.Transaction{
// 		TransactionID:    uuid.New(),
// 		AccountID:        toID,
// 		RelatedAccountID: &fromID,
// 		Amount:           amount,
// 		Type:             model.Transfer,
// 		Note:             "Transfer from account",
// 	}
// 	return repo.Create(ctx, txTo)
// }

// // GetTransactionsByAccount returns paginated results
// func (s *TransactionService) GetTransactionsByAccount(ctx context.Context, accountID uuid.UUID, page, limit int) (map[string]interface{}, error) {
// 	filters := map[string]interface{}{"account_id = ?": accountID}

// 	transactions, err := s.repo.List(ctx, limit, (page-1)*limit, filters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := s.repo.Count(ctx, filters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return utils.PaginatedResponse(transactions, total, limit, page), nil
// }

// // GetAllTransactions returns paginated results
// func (s *TransactionService) GetAllTransactions(ctx context.Context, page, limit int) (map[string]interface{}, error) {
// 	transactions, err := s.repo.List(ctx, limit, (page-1)*limit, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := s.repo.Count(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return utils.PaginatedResponse(transactions, total, limit, page), nil
// }

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

// GetTransactions returns filtered transactions and total count
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

// GetTransactionByID returns only error if not found
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

// service/transactions.go
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
		AccountID:     fromAccID, // record from source account
		Type:          "transfer",
		Amount:        amount,
		CreatedAt:     time.Now(),
	}
	if tx != nil {
		if err := tx.Create(t).Error; err != nil {
			return err
		}
		// Optionally record to destination account as a separate transaction
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
