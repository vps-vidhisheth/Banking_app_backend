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

func InitDB() {
	once.Do(func() {
		dsn := getDSN()

		var err error
		dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf(" Failed to connect to database: %v", err)
		}

		err = dbInstance.AutoMigrate(
			&model.Bank{},
			&model.Customer{},
			&model.Account{},
			&model.Ledger{},
			&model.Transaction{},
		)
		if err != nil {
			log.Fatalf("Failed to migrate models: %v", err)
		}

		log.Println(" Database connected and models migrated successfully.")

		if err := SeedData(dbInstance); err != nil {
			log.Fatalf(" Failed to seed data: %v", err)
		}
	})
}

func GetDB() *gorm.DB {
	if dbInstance == nil {
		InitDB()
	}
	return dbInstance
}

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

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
}
