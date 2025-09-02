package routes

import (
	handler "banking-app/component/ledger/controller"
	service "banking-app/component/ledger/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitLedgerModule(router *mux.Router, db *gorm.DB) {
	// service
	ledgerService := service.NewLedgerService(db) // uses db internally

	// handler
	ledgerHandler := handler.NewLedgerHandler(ledgerService)

	// routes
	handler.RegisterLedgerRoutes(router, ledgerHandler)
}
