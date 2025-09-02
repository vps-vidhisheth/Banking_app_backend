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
	// repository
	bankRepo := repository.NewRepository[model.Bank](db)

	// service
	bankSvc := service.NewBankService(bankRepo, db) // ✅ pass db as second argument

	// handler
	bankHandler := handler.NewBankHandler(bankSvc)

	// routes
	handler.RegisterBankRoutes(router, bankHandler)
}
