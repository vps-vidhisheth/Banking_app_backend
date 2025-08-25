package db

import (
	"banking-app/model"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(database *gorm.DB) error {
	var count int64
	database.Model(&model.Customer{}).Where("role = ?", "admin").Count(&count)

	if count > 0 {
		log.Println(" Admin user already exists, skipping seeding.")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	admin := model.Customer{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@bank.com",
		Password:  string(hashedPassword),
		Role:      "admin",
		IsActive:  true,
	}

	if err := database.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Println(" First admin user created successfully.")
	return nil
}
