package routes

import (
	handler "banking-app/component/customer/controller"
	service "banking-app/component/customer/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitCustomerModule(router *mux.Router, db *gorm.DB) {

	customerService := service.NewCustomerService(db)

	customerHandler := handler.NewCustomerHandler(customerService)

	handler.RegisterCustomerRoutes(router, customerHandler)
}
