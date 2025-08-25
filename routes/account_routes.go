package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// InitAccountModule sets up the account repository, service, handler, and routes
func InitAccountModule(router *mux.Router, db *gorm.DB) {
	// Initialize repositories
	accountRepo := repository.NewAccountRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize dependent services
	ledgerService := service.NewLedgerService(ledgerRepo)
	transactionService := service.NewTransactionService(transactionRepo)

	// Initialize account service
	accountService := service.NewAccountService(accountRepo, ledgerService, transactionService)

	// Initialize handler
	accountHandler := handler.NewAccountHandler(accountService)

	// Register all account routes
	handler.RegisterAccountRoutes(router, accountHandler)
}
