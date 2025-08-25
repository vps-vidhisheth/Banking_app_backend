package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitTransactionModule(router *mux.Router, db *gorm.DB) {
	transactionRepo := repository.NewTransactionRepository(db)

	transactionService := service.NewTransactionService(transactionRepo)

	transactionHandler := handler.NewTransactionHandler(transactionService)

	handler.RegisterTransactionRoutes(router, transactionHandler)
}
