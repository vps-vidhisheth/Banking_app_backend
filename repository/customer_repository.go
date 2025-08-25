package repository

import (
	"banking-app/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(c *model.Customer) error {
	return r.db.Create(c).Error
}

func (r *CustomerRepository) GetByID(id uuid.UUID) (*model.Customer, error) {
	var c model.Customer
	if err := r.db.First(&c, "customer_id = ? AND is_active = ?", id, true).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) Update(c *model.Customer) error {
	return r.db.Save(c).Error
}

func (r *CustomerRepository) Delete(c *model.Customer) error {
	c.IsActive = false
	if err := r.db.Save(c).Error; err != nil {
		return err
	}

	// Soft-delete accounts
	return r.db.Model(&model.Account{}).
		Where("customer_id = ?", c.CustomerID).
		Update("is_active", false).Error
}

func (r *CustomerRepository) GetAll() ([]model.Customer, error) {
	var customers []model.Customer
	if err := r.db.Preload("Accounts").Where("is_active = ?", true).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}
