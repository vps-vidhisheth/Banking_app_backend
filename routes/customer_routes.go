package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitCustomerModule(router *mux.Router, db *gorm.DB) {
	// Initialize repository
	customerRepo := repository.NewCustomerRepository(db)

	// Initialize service
	customerService := service.NewCustomerService(customerRepo)

	// Initialize handler
	customerHandler := handler.NewCustomerHandler(customerService)

	// Register routes
	handler.RegisterCustomerRoutes(router, customerHandler)
}
