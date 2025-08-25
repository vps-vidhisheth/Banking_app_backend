package utils

import (
	"errors"
	"strings"
)

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	return nil
}

func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}
