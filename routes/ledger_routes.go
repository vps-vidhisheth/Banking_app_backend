package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// InitLedgerModule sets up the ledger repository, service, handler, and routes
func InitLedgerModule(router *mux.Router, db *gorm.DB) {
	// Initialize ledger repository
	ledgerRepo := repository.NewLedgerRepository(db)

	// Initialize ledger service with repository
	ledgerService := service.NewLedgerService(ledgerRepo)

	// Initialize handler
	ledgerHandler := handler.NewLedgerHandler(ledgerService)

	// Register all ledger routes
	handler.RegisterLedgerRoutes(router, ledgerHandler)
}
