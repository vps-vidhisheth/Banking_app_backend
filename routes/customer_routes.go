package routes

import (
	handler "banking-app/component/customer/controller"
	service "banking-app/component/customer/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitCustomerModule(router *mux.Router, db *gorm.DB) {
	// service
	customerService := service.NewCustomerService(db) // ✅ only db

	// handler
	customerHandler := handler.NewCustomerHandler(customerService)

	// routes
	handler.RegisterCustomerRoutes(router, customerHandler)
}
