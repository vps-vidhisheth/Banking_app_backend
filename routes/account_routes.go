package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitAccountModule(router *mux.Router, db *gorm.DB) {

	accountRepo := repository.NewAccountRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	ledgerService := service.NewLedgerService(ledgerRepo)
	transactionService := service.NewTransactionService(transactionRepo)

	accountService := service.NewAccountService(accountRepo, ledgerService, transactionService)

	accountHandler := handler.NewAccountHandler(accountService)

	handler.RegisterAccountRoutes(router, accountHandler)
}
