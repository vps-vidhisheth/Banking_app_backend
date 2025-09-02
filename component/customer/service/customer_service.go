// package service

// import (
// 	"banking-app/model"
// 	"banking-app/repository"
// 	"banking-app/utils"
// 	"context"
// 	"errors"
// 	"fmt"
// 	"strings"

// 	"github.com/google/uuid"
// 	"golang.org/x/crypto/bcrypt"
// 	"gorm.io/gorm"
// )

// type CustomerService struct {
// 	repo *repository.Repository[model.Customer]
// 	db   *gorm.DB
// }

// func NewCustomerService(db *gorm.DB) *CustomerService {
// 	repo := repository.NewRepository[model.Customer](db)
// 	cs := &CustomerService{repo: repo, db: db}

// 	// Bootstrap default admin if not present
// 	customers, _ := cs.repo.List(context.Background(), 0, 0, nil)
// 	adminExists := false
// 	for _, c := range customers {
// 		if c.Role == "admin" {
// 			adminExists = true
// 			break
// 		}
// 	}
// 	if !adminExists {
// 		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
// 		_ = cs.repo.Create(context.Background(), &model.Customer{
// 			CustomerID: uuid.New(),
// 			FirstName:  "Admin",
// 			LastName:   "User",
// 			Email:      "admin@bank.com",
// 			Password:   string(hashed),
// 			Role:       "admin",
// 			IsActive:   true,
// 		})
// 	}

// 	return cs
// }

// // ---------------- AUTHENTICATION ----------------
// func (cs *CustomerService) Authenticate(ctx context.Context, email, password string) (*model.Customer, error) {
// 	customers, err := cs.repo.List(ctx, 0, 0, nil)
// 	if err != nil {
// 		return nil, errors.New("failed to fetch customers")
// 	}

// 	email = strings.TrimSpace(strings.ToLower(email))
// 	for i := range customers {
// 		c := &customers[i]
// 		if c.IsActive && strings.EqualFold(c.Email, email) {
// 			if err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)); err != nil {
// 				return nil, errors.New("invalid email or password")
// 			}
// 			return c, nil
// 		}
// 	}
// 	return nil, errors.New("invalid email or password")
// }

// // ---------------- CREATE CUSTOMER WITH TRANSACTION ----------------
// func (cs *CustomerService) CreateCustomer(ctx context.Context, firstName, lastName, email, password, role string) (*model.Customer, error) {
// 	email = strings.TrimSpace(strings.ToLower(email))

// 	returnValue := &model.Customer{}
// 	err := cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		txRepo := repository.NewRepository[model.Customer](tx)

// 		existing, _ := txRepo.List(ctx, 0, 0, map[string]interface{}{"email = ?": email})
// 		if len(existing) > 0 {
// 			return errors.New("customer with this email already exists")
// 		}

// 		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 		if err != nil {
// 			return fmt.Errorf("failed to hash password: %v", err)
// 		}

// 		c := &model.Customer{
// 			CustomerID: uuid.New(),
// 			FirstName:  firstName,
// 			LastName:   lastName,
// 			Email:      email,
// 			Password:   string(hashed),
// 			Role:       role,
// 			IsActive:   true,
// 		}

// 		if err := txRepo.Create(ctx, c); err != nil {
// 			return err
// 		}

// 		*returnValue = *c
// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}
// 	return returnValue, nil
// }

// // ---------------- GET CUSTOMER BY ID ----------------
// func (cs *CustomerService) GetCustomerByID(ctx context.Context, id uuid.UUID) (*model.Customer, error) {
// 	c, err := cs.repo.GetOne(ctx, "customer_id = ?", id)
// 	if err != nil {
// 		return nil, fmt.Errorf("customer not found: %v", err)
// 	}
// 	return c, nil
// }

// // ---------------- UPDATE CUSTOMER WITH TRANSACTION ----------------
// func (cs *CustomerService) UpdateCustomer(ctx context.Context, id uuid.UUID, firstName, lastName, email, role string, isActive *bool) error {
// 	return cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		txRepo := repository.NewRepository[model.Customer](tx)

// 		c, err := txRepo.GetOne(ctx, "customer_id = ?", id)
// 		if err != nil || c == nil {
// 			return errors.New("customer not found")
// 		}

// 		if firstName != "" {
// 			c.FirstName = firstName
// 		}
// 		if lastName != "" {
// 			c.LastName = lastName
// 		}
// 		if email != "" {
// 			c.Email = strings.TrimSpace(strings.ToLower(email))
// 		}
// 		if role != "" {
// 			c.Role = role
// 		}
// 		if isActive != nil {
// 			c.IsActive = *isActive
// 		}

// 		return txRepo.Update(ctx, c)
// 	})
// }

// // ---------------- SOFT DELETE CUSTOMER ----------------
// func (cs *CustomerService) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
// 	return cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		txRepo := repository.NewRepository[model.Customer](tx)
// 		return txRepo.Delete(ctx, "customer_id = ?", id)
// 	})
// }

// // ---------------- GET ALL CUSTOMERS ----------------
// func (cs *CustomerService) GetAllCustomers(ctx context.Context) ([]model.Customer, error) {
// 	return cs.repo.List(ctx, 0, 0, nil)
// }

// // ---------------- GET ALL CUSTOMERS WITH PAGINATION ----------------
// func (cs *CustomerService) GetAllCustomersPaginated(ctx context.Context, lastName string, page, limit int) (map[string]interface{}, error) {
// 	allCustomers, err := cs.repo.List(ctx, 0, 0, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	filtered := []model.Customer{}
// 	for _, c := range allCustomers {
// 		if !c.IsActive {
// 			continue
// 		}
// 		if lastName == "" || strings.Contains(strings.ToLower(c.LastName), strings.ToLower(lastName)) {
// 			filtered = append(filtered, c)
// 		}
// 	}

// 	total := int64(len(filtered))
// 	start := (page - 1) * limit
// 	if start >= len(filtered) {
// 		start = len(filtered)
// 	}
// 	end := start + limit
// 	if end > len(filtered) {
// 		end = len(filtered)
// 	}
// 	paginated := filtered[start:end]

// 	return utils.PaginatedResponse(paginated, total, limit, page), nil
// }

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

	// Bootstrap default admin if not present
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

// ---------------- AUTHENTICATION ----------------
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

// ---------------- CREATE CUSTOMER ----------------
func (cs *CustomerService) CreateCustomer(ctx context.Context, firstName, lastName, email, password, role string) error {
	email = strings.TrimSpace(strings.ToLower(email))

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

// ---------------- UPDATE CUSTOMER ----------------
func (cs *CustomerService) UpdateCustomer(ctx context.Context, id uuid.UUID, firstName, lastName, email, role string, isActive *bool) error {
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

// ---------------- DELETE CUSTOMER ----------------
func (cs *CustomerService) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	return cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewRepository[model.Customer](tx)
		return txRepo.Delete(ctx, "customer_id = ?", id)
	})
}

// ---------------- CHECK CUSTOMER EXISTS WITH FILTERS ----------------
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

// ---------------- GET CUSTOMER BY EMAIL ----------------
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

// ListCustomers returns paginated customers with optional filters
func (cs *CustomerService) ListCustomers(ctx context.Context, limit, offset int, filters map[string]string) ([]model.Customer, error) {
	queryFilters := make(map[string]interface{})
	for k, v := range filters {
		if v != "" {
			queryFilters[k+" LIKE ?"] = "%" + v + "%"
		}
	}
	return cs.repo.List(ctx, limit, offset, queryFilters)
}

// CountCustomers returns total count of customers matching filters
func (cs *CustomerService) CountCustomers(ctx context.Context, filters map[string]string) (int64, error) {
	queryFilters := make(map[string]interface{})
	for k, v := range filters {
		if v != "" {
			queryFilters[k+" LIKE ?"] = "%" + v + "%"
		}
	}
	return cs.repo.Count(ctx, queryFilters)
}

// GetCustomerByID returns the customer object by ID
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
