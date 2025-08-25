package service

import (
	"banking-app/middleware"
	"banking-app/model"

	"github.com/google/uuid"
)

type AuthService struct {
	CustomerService *CustomerService
}

func NewAuthService(cs *CustomerService) *AuthService {
	return &AuthService{
		CustomerService: cs,
	}
}

func (a *AuthService) GenerateToken(customer *model.Customer) (string, error) {
	userID, err := uuid.Parse(customer.CustomerID.String())
	if err != nil {
		return "", err
	}
	return middleware.GenerateToken(userID, customer.Role)
}

func (a *AuthService) Authenticate(email, password string) (string, error) {
	customer, err := a.CustomerService.Authenticate(email, password)
	if err != nil {
		return "", err
	}
	return a.GenerateToken(customer)
}
