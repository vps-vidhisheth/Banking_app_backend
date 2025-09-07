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

func NewAuthService(cs *customerSvc.CustomerService) *AuthService {
	return &AuthService{
		CustomerService: cs,
	}
}

func (a *AuthService) Authenticate(ctx context.Context, email, password string) error {
	if err := a.CustomerService.Authenticate(ctx, email, password); err != nil {
		return err
	}
	return nil
}

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
