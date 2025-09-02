package routes

import (
	accountHandler "banking-app/component/account/controller"
	accountService "banking-app/component/account/service"
	ledgerService "banking-app/component/ledger/service"
	transactionService "banking-app/component/transactions/service" // ✅ plural
	"banking-app/repository"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitAccountModule(router *mux.Router, db *gorm.DB) {
	// Create UnitOfWork from DB
	uow := repository.NewUnitOfWork(db)

	// Services
	ledgerSvc := ledgerService.NewLedgerService(db)
	transactionSvc := transactionService.NewTransactionService(db)

	accountSvc := accountService.NewAccountService(uow, ledgerSvc, transactionSvc)

	// Handler
	accHandler := accountHandler.NewAccountHandler(accountSvc)

	// Routes
	accountHandler.RegisterAccountRoutes(router, accHandler)
}
