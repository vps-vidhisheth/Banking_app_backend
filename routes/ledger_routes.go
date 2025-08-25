package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitLedgerModule(router *mux.Router, db *gorm.DB) {

	ledgerRepo := repository.NewLedgerRepository(db)

	ledgerService := service.NewLedgerService(ledgerRepo)

	ledgerHandler := handler.NewLedgerHandler(ledgerService)

	handler.RegisterLedgerRoutes(router, ledgerHandler)
}
