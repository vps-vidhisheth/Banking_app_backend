package routes

import (
	handler "banking-app/component/transactions/controller"
	service "banking-app/component/transactions/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitTransactionModule(router *mux.Router, db *gorm.DB) {

	transactionService := service.NewTransactionService(db)

	transactionHandler := handler.NewTransactionHandler(transactionService)

	handler.RegisterTransactionRoutes(router, transactionHandler)
}
