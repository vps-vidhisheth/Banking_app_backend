package routes

import (
	handler "banking-app/component/ledger/controller"
	service "banking-app/component/ledger/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitLedgerModule(router *mux.Router, db *gorm.DB) {

	ledgerService := service.NewLedgerService(db)

	ledgerHandler := handler.NewLedgerHandler(ledgerService)

	handler.RegisterLedgerRoutes(router, ledgerHandler)
}
