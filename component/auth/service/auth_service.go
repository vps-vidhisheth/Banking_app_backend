// package service

// import (
// 	customerSvc "banking-app/component/customer/service"

// 	"banking-app/middleware"
// 	"banking-app/model"

// 	"context" // << add this

// 	"github.com/google/uuid"
// )

// type AuthService struct {
// 	CustomerService *customerSvc.CustomerService
// }

// // constructor
// func NewAuthService(cs *customerSvc.CustomerService) *AuthService {
// 	return &AuthService{
// 		CustomerService: cs,
// 	}
// }

// // GenerateToken generates a JWT token for the authenticated customer
// func (a *AuthService) GenerateToken(customer *model.Customer) (string, error) {
// 	userID, err := uuid.Parse(customer.CustomerID.String())
// 	if err != nil {
// 		return "", err
// 	}
// 	return middleware.GenerateToken(userID, customer.Role)
// }

// // Authenticate validates credentials and returns a JWT token
// func (a *AuthService) Authenticate(ctx context.Context, email, password string) (string, error) {
// 	customer, err := a.CustomerService.Authenticate(ctx, email, password)
// 	if err != nil {
// 		return "", err
// 	}
// 	return a.GenerateToken(customer)
// }

package service

import (
	customerSvc "banking-app/component/customer/service"
	"banking-app/middleware"
	"banking-app/model"
	"context"
	"errors"

	"github.com/google/uuid"
)

type AuthService struct {
	CustomerService *customerSvc.CustomerService
}

// constructor
func NewAuthService(cs *customerSvc.CustomerService) *AuthService {
	return &AuthService{
		CustomerService: cs,
	}
}

// Authenticate validates credentials and returns error only
func (a *AuthService) Authenticate(ctx context.Context, email, password string) error {
	// Use CustomerService to check credentials
	if err := a.CustomerService.Authenticate(ctx, email, password); err != nil {
		return err
	}
	return nil
}

// GenerateToken returns error if customer is invalid
func (a *AuthService) GenerateToken(customer *model.Customer) (string, error) {
	if customer == nil {
		return "", errors.New("invalid customer")
	}

	userID, err := uuid.Parse(customer.CustomerID.String())
	if err != nil {
		return "", err
	}
	token, err := middleware.GenerateToken(userID, customer.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}
