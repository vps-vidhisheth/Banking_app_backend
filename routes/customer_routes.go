package routes

import (
	"banking-app/handler"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitCustomerModule(router *mux.Router, db *gorm.DB) {
	customerRepo := repository.NewCustomerRepository(db)

	customerService := service.NewCustomerService(customerRepo)

	customerHandler := handler.NewCustomerHandler(customerService)

	handler.RegisterCustomerRoutes(router, customerHandler)
}
