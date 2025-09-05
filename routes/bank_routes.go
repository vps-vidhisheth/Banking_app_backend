package routes

import (
	handler "banking-app/component/banks/controller"
	service "banking-app/component/banks/service"
	"banking-app/model"
	"banking-app/repository"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitBankModule(router *mux.Router, db *gorm.DB) {

	bankRepo := repository.NewRepository[model.Bank](db)

	bankSvc := service.NewBankService(bankRepo, db)

	bankHandler := handler.NewBankHandler(bankSvc)

	handler.RegisterBankRoutes(router, bankHandler)
}
