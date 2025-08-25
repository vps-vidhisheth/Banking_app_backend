package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// InitTransactionModule sets up the transaction repository, service, handler, and routes
func InitTransactionModule(router *mux.Router, db *gorm.DB) {
	// Initialize transaction repository
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize transaction service with repository
	transactionService := service.NewTransactionService(transactionRepo)

	// Initialize handler
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Register all transaction routes
	handler.RegisterTransactionRoutes(router, transactionHandler)
}
