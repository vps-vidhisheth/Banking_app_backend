package routes

import (
	handler "banking-app/component/transactions/controller"
	service "banking-app/component/transactions/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitTransactionModule(router *mux.Router, db *gorm.DB) {
	// service (only pass db)
	transactionService := service.NewTransactionService(db)

	// handler
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// routes
	handler.RegisterTransactionRoutes(router, transactionHandler)
}
