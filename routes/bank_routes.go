package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitBankModule(router *mux.Router, db *gorm.DB) {

	bankRepo := repository.NewBankRepository(db)

	bankService := service.NewBankService(bankRepo)

	bankHandler := handler.NewBankHandler(bankService)

	handler.RegisterBankRoutes(router, bankHandler)
}
