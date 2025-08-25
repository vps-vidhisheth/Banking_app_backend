package service

import (
	"banking-app/apperror"
	"banking-app/model"
	"banking-app/repository"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type BankService struct {
	repo *repository.BankRepository
}

func NewBankService(repo *repository.BankRepository) *BankService {
	return &BankService{repo: repo}
}

func (s *BankService) CreateBank(name string) (*model.Bank, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("bank name cannot be empty")
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

	if err := s.repo.Create(bank); err != nil {
		return nil, err
	}
	return bank, nil
}

func (s *BankService) GetBankByID(bankID uuid.UUID) (*model.Bank, error) {
	bank, err := s.repo.GetByID(bankID)
	if err != nil {
		return nil, err
	}
	if bank == nil {
		return nil, apperror.NewNotFoundError("bank not found")
	}
	return bank, nil
}

func (s *BankService) ListBanks() ([]*model.Bank, error) {
	return s.repo.List()
}

func (s *BankService) UpdateBank(bankID uuid.UUID, newName string) error {
	bank, err := s.repo.GetByID(bankID)
	if err != nil {
		return err
	}
	if bank == nil {
		return apperror.NewNotFoundError("bank not found")
	}

	newName = strings.TrimSpace(newName)
	if newName == "" {
		return errors.New("bank name cannot be empty")
	}

	bank.Name = newName
	if len(newName) >= 3 {
		bank.Abbreviation = strings.ToUpper(newName[:3])
	} else {
		bank.Abbreviation = strings.ToUpper(newName)
	}

	return s.repo.Update(bank)
}

func (s *BankService) DeleteBank(bankID uuid.UUID) error {
	bank, err := s.repo.GetByID(bankID)
	if err != nil {
		return err
	}
	if bank == nil {
		return apperror.NewNotFoundError("bank not found")
	}

	for _, acc := range bank.Accounts {
		if acc.IsActive {
			return errors.New("cannot delete bank with active accounts")
		}
	}

	return s.repo.Delete(bankID)
}
