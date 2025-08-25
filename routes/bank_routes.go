package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// InitBankModule sets up the bank repository, service, handler, and routes
func InitBankModule(router *mux.Router, db *gorm.DB) {
	// Initialize repository (database-backed)
	bankRepo := repository.NewBankRepository(db)

	// Initialize bank service
	bankService := service.NewBankService(bankRepo)

	// Initialize handler
	bankHandler := handler.NewBankHandler(bankService)

	// Register all bank routes
	handler.RegisterBankRoutes(router, bankHandler)
}
