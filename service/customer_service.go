package service

import (
	"banking-app/model"
	"banking-app/repository"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	cs := &CustomerService{repo: repo}

	customers, _ := cs.repo.GetAll()
	adminExists := false
	for _, c := range customers {
		if c.Role == "admin" {
			adminExists = true
			break
		}
	}
	if !adminExists {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_ = cs.repo.Create(&model.Customer{
			CustomerID: uuid.New(),
			FirstName:  "Admin",
			LastName:   "User",
			Email:      "admin@bank.com",
			Password:   string(hashed),
			Role:       "admin",
			IsActive:   true,
		})
	}

	return cs
}

func (cs *CustomerService) Authenticate(email, password string) (*model.Customer, error) {
	allCustomers, err := cs.repo.GetAll()
	if err != nil {
		return nil, errors.New("failed to fetch customers")
	}

	email = strings.TrimSpace(strings.ToLower(email))
	for i := range allCustomers {
		c := &allCustomers[i]
		if c.IsActive && strings.EqualFold(c.Email, email) {
			if err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)); err != nil {
				return nil, errors.New("invalid email or password")
			}
			return c, nil
		}
	}

	return nil, errors.New("invalid email or password")
}

func (cs *CustomerService) CreateCustomer(firstName, lastName, email, password, role string) (*model.Customer, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	allCustomers, _ := cs.repo.GetAll()
	for _, c := range allCustomers {
		if c.Email == email {
			return nil, errors.New("customer with this email already exists")
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	c := &model.Customer{
		CustomerID: uuid.New(),
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		Password:   string(hashed),
		Role:       role,
		IsActive:   true,
	}

	if err := cs.repo.Create(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (cs *CustomerService) GetCustomerByID(id uuid.UUID) (*model.Customer, error) {
	c, err := cs.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %v", err)
	}
	return c, nil
}

func (cs *CustomerService) UpdateCustomer(id uuid.UUID, firstName, lastName, email, role string, isActive *bool) error {
	c, err := cs.repo.GetByID(id)
	if err != nil {
		return err
	}

	if firstName != "" {
		c.FirstName = firstName
	}
	if lastName != "" {
		c.LastName = lastName
	}
	if email != "" {
		c.Email = strings.TrimSpace(strings.ToLower(email))
	}
	if role != "" {
		c.Role = role
	}
	if isActive != nil {
		c.IsActive = *isActive
	}

	return cs.repo.Update(c)
}

func (cs *CustomerService) DeleteCustomer(id uuid.UUID) error {
	c, err := cs.repo.GetByID(id)
	if err != nil {
		return err
	}
	return cs.repo.Delete(c)
}

func (cs *CustomerService) GetAllCustomers() ([]model.Customer, error) {
	return cs.repo.GetAll()
}

func (cs *CustomerService) GetAllCustomersPaginated(lastName string, page, limit int) ([]model.Customer, int64, error) {
	allCustomers, err := cs.repo.GetAll()
	if err != nil {
		return nil, 0, err
	}

	filtered := []model.Customer{}
	for _, c := range allCustomers {
		if !c.IsActive {
			continue
		}
		if lastName == "" || strings.Contains(strings.ToLower(c.LastName), strings.ToLower(lastName)) {
			filtered = append(filtered, c)
		}
	}

	total := int64(len(filtered))

	start := (page - 1) * limit
	if start >= len(filtered) {
		return []model.Customer{}, total, nil
	}

	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	paginated := filtered[start:end]
	return paginated, total, nil
}
