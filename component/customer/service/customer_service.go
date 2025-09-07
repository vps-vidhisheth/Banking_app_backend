package service

import (
	"banking-app/model"
	"banking-app/repository"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomerService struct {
	repo *repository.Repository[model.Customer]
	db   *gorm.DB
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	repo := repository.NewRepository[model.Customer](db)
	cs := &CustomerService{repo: repo, db: db}

	customers, _ := cs.repo.List(context.Background(), 0, 0, nil)
	adminExists := false
	for _, c := range customers {
		if c.Role == "admin" {
			adminExists = true
			break
		}
	}
	if !adminExists {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_ = cs.repo.Create(context.Background(), &model.Customer{
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

func (cs *CustomerService) Authenticate(ctx context.Context, email, password string) error {
	email = strings.TrimSpace(strings.ToLower(email))

	customers, err := cs.repo.List(ctx, 0, 0, nil)
	if err != nil {
		return errors.New("failed to fetch customers")
	}

	for i := range customers {
		c := &customers[i]
		if c.IsActive && strings.EqualFold(c.Email, email) {
			if err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)); err != nil {
				return errors.New("invalid email or password")
			}
			return nil
		}
	}

	return errors.New("invalid email or password")
}

func (cs *CustomerService) CreateCustomer(ctx context.Context, firstName, lastName, email, password, role string) error {
	email = strings.TrimSpace(strings.ToLower(email))

	if !strings.HasSuffix(email, ".com") {
		return errors.New("email must end with .com")
	}

	return cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewRepository[model.Customer](tx)

		existing, _ := txRepo.List(ctx, 0, 0, map[string]interface{}{"email = ?": email})
		if len(existing) > 0 {
			return errors.New("customer with this email already exists")
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
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

		return txRepo.Create(ctx, c)
	})
}

func (cs *CustomerService) UpdateCustomer(ctx context.Context, id uuid.UUID, firstName, lastName, email, role string, isActive *bool) error {

	if email != "" && !strings.HasSuffix(strings.TrimSpace(strings.ToLower(email)), ".com") {
		return errors.New("email must end with .com")
	}

	return cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewRepository[model.Customer](tx)

		c, err := txRepo.GetOne(ctx, "customer_id = ?", id)
		if err != nil || c == nil {
			return errors.New("customer not found")
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

		return txRepo.Update(ctx, c)
	})
}

func (cs *CustomerService) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	return cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewRepository[model.Customer](tx)
		return txRepo.Delete(ctx, "customer_id = ?", id)
	})
}

func (cs *CustomerService) CheckCustomerExists(ctx context.Context, filters map[string]string) error {
	if filters == nil {
		return errors.New("no filters provided")
	}

	queryFilters := make(map[string]interface{})
	for k, v := range filters {
		if v != "" {
			queryFilters[k+" = ?"] = v
		}
	}

	customers, err := cs.repo.List(ctx, 0, 0, queryFilters)
	if err != nil {
		return err
	}
	if len(customers) == 0 {
		return errors.New("no matching customer found")
	}
	return nil
}

func (cs *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*model.Customer, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	customers, err := cs.repo.List(ctx, 0, 0, map[string]interface{}{"email = ?": email})
	if err != nil {
		return nil, err
	}
	if len(customers) == 0 {
		return nil, errors.New("customer not found")
	}
	return &customers[0], nil
}

func (cs *CustomerService) ListCustomers(ctx context.Context, limit, offset int, filters map[string]string) ([]model.Customer, error) {
	query := cs.db.WithContext(ctx).Model(&model.Customer{})

	if search, ok := filters["search"]; ok && search != "" {
		like := "%" + search + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", like, like, like)
	}

	if v := filters["first_name"]; v != "" {
		query = query.Where("first_name LIKE ?", "%"+v+"%")
	}
	if v := filters["last_name"]; v != "" {
		query = query.Where("last_name LIKE ?", "%"+v+"%")
	}
	if v := filters["email"]; v != "" {
		query = query.Where("email LIKE ?", "%"+v+"%")
	}
	if v := filters["role"]; v != "" {
		query = query.Where("role = ?", v)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	query = query.Order("customer_id ASC")

	var customers []model.Customer
	if err := query.Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (cs *CustomerService) CountCustomers(ctx context.Context, filters map[string]string) (int64, error) {
	query := cs.db.WithContext(ctx).Model(&model.Customer{})

	if search, ok := filters["search"]; ok && search != "" {
		like := "%" + search + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", like, like, like)
	}

	if v := filters["first_name"]; v != "" {
		query = query.Where("first_name LIKE ?", "%"+v+"%")
	}
	if v := filters["last_name"]; v != "" {
		query = query.Where("last_name LIKE ?", "%"+v+"%")
	}
	if v := filters["email"]; v != "" {
		query = query.Where("email LIKE ?", "%"+v+"%")
	}
	if v := filters["role"]; v != "" {
		query = query.Where("role = ?", v)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (cs *CustomerService) GetCustomerByID(ctx context.Context, id uuid.UUID) (*model.Customer, error) {
	customer, err := cs.repo.GetOne(ctx, "customer_id = ?", id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("customer not found")
	}
	return customer, nil
}
