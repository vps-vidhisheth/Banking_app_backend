// package service

// import (
// 	"banking-app/apperror"
// 	"banking-app/model"
// 	"banking-app/repository"
// 	"banking-app/utils"
// 	"context"
// 	"errors"
// 	"strings"
// 	"time"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type BankService struct {
// 	repo *repository.Repository[model.Bank]
// 	db   *gorm.DB
// }

// func NewBankService(repo *repository.Repository[model.Bank], db *gorm.DB) *BankService {
// 	return &BankService{repo: repo, db: db}
// }

// // ----------------- Create Bank -----------------
// func (s *BankService) CreateBank(ctx context.Context, name string) (*model.Bank, error) {
// 	name = strings.TrimSpace(name)
// 	if name == "" {
// 		return nil, errors.New("bank name cannot be empty")
// 	}

// 	bank := &model.Bank{
// 		BankID:   uuid.New(),
// 		Name:     name,
// 		IsActive: true,
// 	}

// 	if len(name) >= 3 {
// 		bank.Abbreviation = strings.ToUpper(name[:3])
// 	} else {
// 		bank.Abbreviation = strings.ToUpper(name)
// 	}

// 	// Transaction ensures rollback on failure
// 	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		tempRepo := repository.NewRepository[model.Bank](tx)
// 		if err := tempRepo.Create(ctx, bank); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bank, nil
// }

// // ----------------- Get Bank -----------------
// func (s *BankService) GetBankByID(ctx context.Context, bankID uuid.UUID) (*model.Bank, error) {
// 	bank, err := s.repo.GetOne(ctx, "bank_id = ?", bankID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if bank == nil {
// 		return nil, apperror.NewNotFoundError("bank not found")
// 	}
// 	return bank, nil
// }

// // ----------------- List Banks with Pagination -----------------
// func (s *BankService) ListBanks(ctx context.Context, limit, offset int) (map[string]interface{}, error) {
// 	filters := map[string]interface{}{"is_active = ?": true}

// 	banks, err := s.repo.List(ctx, limit, offset, filters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := s.repo.Count(ctx, filters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return utils.PaginatedResponse(banks, total, limit, offset), nil
// }

// // ----------------- Update Bank -----------------
// func (s *BankService) UpdateBank(ctx context.Context, bankID uuid.UUID, newName string) error {
// 	newName = strings.TrimSpace(newName)
// 	if newName == "" {
// 		return errors.New("bank name cannot be empty")
// 	}

// 	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		tempRepo := repository.NewRepository[model.Bank](tx)
// 		bank, err := tempRepo.GetOne(ctx, "bank_id = ?", bankID)
// 		if err != nil {
// 			return err
// 		}
// 		if bank == nil {
// 			return apperror.NewNotFoundError("bank not found")
// 		}

// 		bank.Name = newName
// 		if len(newName) >= 3 {
// 			bank.Abbreviation = strings.ToUpper(newName[:3])
// 		} else {
// 			bank.Abbreviation = strings.ToUpper(newName)
// 		}

// 		if err := tempRepo.Update(ctx, bank); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

// // ----------------- Delete Bank (soft-delete check for active accounts) -----------------
// func (s *BankService) DeleteBank(ctx context.Context, bankID uuid.UUID) error {
// 	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		bankRepo := repository.NewRepository[model.Bank](tx)
// 		accountRepo := repository.NewRepository[model.Account](tx)

// 		// Fetch bank
// 		bank, err := bankRepo.GetOne(ctx, "bank_id = ?", bankID)
// 		if err != nil {
// 			return err
// 		}
// 		if bank == nil {
// 			return apperror.NewNotFoundError("bank not found")
// 		}

// 		// Soft delete all active accounts
// 		accounts, err := accountRepo.List(ctx, 0, 0, map[string]interface{}{
// 			"bank_id":   bankID,
// 			"is_active": true,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		for _, acc := range accounts {
// 			acc.IsActive = false
// 			acc.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
// 			if err := accountRepo.Update(ctx, &acc); err != nil {
// 				return err
// 			}
// 		}

// 		// Soft delete the bank
// 		bank.IsActive = false
// 		bank.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
// 		if err := bankRepo.Update(ctx, bank); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
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

type BankService struct {
	repo *repository.Repository[model.Bank]
	db   *gorm.DB
}

func NewBankService(repo *repository.Repository[model.Bank], db *gorm.DB) *BankService {
	return &BankService{repo: repo, db: db}
}

// ----------------- Create Bank -----------------
func (s *BankService) CreateBank(ctx context.Context, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("bank name cannot be empty")
	}

	bank := &model.Bank{
		BankID:   uuid.New(),
		Name:     name,
		IsActive: true,
	}

	if len(name) >= 3 {
		bank.Abbreviation = strings.ToUpper(name[:3])
	} else {
		bank.Abbreviation = strings.ToUpper(name)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tempRepo := repository.NewRepository[model.Bank](tx)
		return tempRepo.Create(ctx, bank)
	})
}

// ----------------- Update Bank -----------------
func (s *BankService) UpdateBank(ctx context.Context, bankID uuid.UUID, newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return errors.New("bank name cannot be empty")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tempRepo := repository.NewRepository[model.Bank](tx)
		bank, err := tempRepo.GetOne(ctx, "bank_id = ?", bankID)
		if err != nil {
			return err
		}
		if bank == nil {
			return errors.New("bank not found")
		}

		bank.Name = newName
		if len(newName) >= 3 {
			bank.Abbreviation = strings.ToUpper(newName[:3])
		} else {
			bank.Abbreviation = strings.ToUpper(newName)
		}

		return tempRepo.Update(ctx, bank)
	})
}

// ----------------- Delete Bank (soft-delete active accounts first) -----------------
func (s *BankService) DeleteBank(ctx context.Context, bankID uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		bankRepo := repository.NewRepository[model.Bank](tx)
		accountRepo := repository.NewRepository[model.Account](tx)

		bank, err := bankRepo.GetOne(ctx, "bank_id = ?", bankID)
		if err != nil {
			return err
		}
		if bank == nil {
			return errors.New("bank not found")
		}

		// Soft delete all active accounts
		accounts, err := accountRepo.List(ctx, 0, 0, map[string]interface{}{
			"bank_id":   bankID,
			"is_active": true,
		})
		if err != nil {
			return err
		}
		for _, acc := range accounts {
			acc.IsActive = false
			acc.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			if err := accountRepo.Update(ctx, &acc); err != nil {
				return err
			}
		}

		// Soft delete the bank
		bank.IsActive = false
		bank.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
		return bankRepo.Update(ctx, bank)
	})
}

// ----------------- Check Bank Exists -----------------
func (s *BankService) CheckBankExists(ctx context.Context, bankID uuid.UUID) error {
	bank, err := s.repo.GetOne(ctx, "bank_id = ? AND is_active = ?", bankID, true)
	if err != nil {
		return err
	}
	if bank == nil {
		return errors.New("bank not found")
	}
	return nil
}

// ----------------- Check Any Active Banks -----------------
func (s *BankService) CheckAnyBanks(ctx context.Context) error {
	filters := map[string]interface{}{"is_active = ?": true}
	banks, err := s.repo.List(ctx, 1, 0, filters)
	if err != nil {
		return err
	}
	if len(banks) == 0 {
		return errors.New("no active banks found")
	}
	return nil
}

// ----------------- Check Banks with Filters -----------------
// Supports filters: name, abbreviation, is_active
func (s *BankService) CheckBanksWithFilters(ctx context.Context, filters map[string]string) error {
	query := make(map[string]interface{})
	for key, value := range filters {
		if value == "" {
			continue
		}
		switch key {
		case "name":
			query["name LIKE ?"] = "%" + strings.TrimSpace(value) + "%"
		case "abbreviation":
			query["abbreviation LIKE ?"] = "%" + strings.TrimSpace(value) + "%"
		case "is_active":
			if strings.ToLower(value) == "true" {
				query["is_active = ?"] = true
			} else if strings.ToLower(value) == "false" {
				query["is_active = ?"] = false
			}
		}
	}

	banks, err := s.repo.List(ctx, 1, 0, query)
	if err != nil {
		return err
	}
	if len(banks) == 0 {
		return errors.New("no banks found with given filters")
	}
	return nil
}

// ----------------- List Active Banks with Pagination -----------------
func (s *BankService) ListBanksPaginated(ctx context.Context, limit, offset int) ([]model.Bank, int64, error) {
	filters := map[string]interface{}{"is_active = ?": true}

	banks, err := s.repo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return banks, total, nil
}

func (s *BankService) ListBanksWithFilters(ctx context.Context, limit, offset int, filters map[string]string) ([]model.Bank, int64, error) {
	query := make(map[string]interface{})
	for key, value := range filters {
		if value == "" {
			continue
		}
		switch key {
		case "name":
			query["name LIKE ?"] = "%" + strings.TrimSpace(value) + "%"
		case "abbreviation":
			query["abbreviation LIKE ?"] = "%" + strings.TrimSpace(value) + "%"
		case "is_active":
			if strings.ToLower(value) == "true" {
				query["is_active = ?"] = true
			} else if strings.ToLower(value) == "false" {
				query["is_active = ?"] = false
			}
		}
	}

	banks, err := s.repo.List(ctx, limit, offset, query)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return banks, total, nil
}
