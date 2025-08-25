package db

import (
	"banking-app/model"
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// InitDB initializes the database connection (singleton)
func InitDB() {
	once.Do(func() {
		dsn := getDSN()

		var err error
		dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf(" Failed to connect to database: %v", err)
		}

		// AutoMigrate all models
		err = dbInstance.AutoMigrate(
			&model.Bank{},
			&model.Customer{},
			&model.Account{},
			&model.Ledger{},
			&model.Transaction{},
		)
		if err != nil {
			log.Fatalf("❌ Failed to migrate models: %v", err)
		}

		log.Println("✅ Database connected and models migrated successfully.")

		// Seed initial data if needed
		if err := SeedData(dbInstance); err != nil {
			log.Fatalf("❌ Failed to seed data: %v", err)
		}
	})
}

// GetDB returns the DB instance
func GetDB() *gorm.DB {
	if dbInstance == nil {
		InitDB()
	}
	return dbInstance
}

// getDSN builds the Data Source Name from env variables or defaults
func getDSN() string {
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = "pass@123"
	}
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "banking_app_db"
	}

	// Return DSN in correct format for MySQL + GORM
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
}
