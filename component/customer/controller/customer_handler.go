// package handler

// import (
// 	"banking-app/component/customer/service"
// 	"banking-app/middleware"
// 	"banking-app/utils"
// 	"banking-app/web"
// 	"net/http"

// 	"github.com/google/uuid"
// 	"github.com/gorilla/mux"
// )

// type CustomerHandler struct {
// 	Service *service.CustomerService
// }

// func NewCustomerHandler(cs *service.CustomerService) *CustomerHandler {
// 	return &CustomerHandler{Service: cs}
// }

// // adminOnly helper
// func (h *CustomerHandler) adminOnly(w http.ResponseWriter, r *http.Request) bool {
// 	claims, ok := middleware.GetUserClaims(r)
// 	if !ok || claims.Role != "admin" {
// 		web.RespondErrorMessage(w, http.StatusForbidden, "only admin can perform this action")
// 		return false
// 	}
// 	return true
// }

// func (h *CustomerHandler) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.adminOnly(w, r) {
// 		return
// 	}

// 	var req struct {
// 		FirstName string `json:"first_name"`
// 		LastName  string `json:"last_name"`
// 		Email     string `json:"email"`
// 		Password  string `json:"password"`
// 		Role      string `json:"role"`
// 		IsActive  *bool  `json:"is_active,omitempty"`
// 	}

// 	if err := web.UnmarshalJSON(r, &req); err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}

// 	if req.Role == "" {
// 		req.Role = "customer"
// 	}

// 	ctx := r.Context()
// 	cust, err := h.Service.CreateCustomer(ctx, req.FirstName, req.LastName, req.Email, req.Password, req.Role)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if req.IsActive != nil {
// 		_ = h.Service.UpdateCustomer(ctx, cust.CustomerID, "", "", "", "", req.IsActive)
// 	}

// 	web.RespondJSON(w, http.StatusCreated, cust)
// }

// // Get all customers with pagination + filtering
// func (h *CustomerHandler) GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
// 	lastName := r.URL.Query().Get("last_name")
// 	pagination := utils.GetPaginationParams(r, 2, 1)

// 	ctx := r.Context()
// 	resp, err := h.Service.GetAllCustomersPaginated(ctx, lastName, pagination.Offset, pagination.Limit)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, resp)
// }

// func (h *CustomerHandler) GetCustomerByIDHandler(w http.ResponseWriter, r *http.Request) {
// 	idStr := mux.Vars(r)["id"]
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer ID")
// 		return
// 	}

// 	ctx := r.Context()
// 	cust, err := h.Service.GetCustomerByID(ctx, id)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, cust)
// }

// func (h *CustomerHandler) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.adminOnly(w, r) {
// 		return
// 	}

// 	idStr := mux.Vars(r)["id"]
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer ID")
// 		return
// 	}

// 	var req struct {
// 		FirstName string `json:"first_name"`
// 		LastName  string `json:"last_name"`
// 		Email     string `json:"email"`
// 		Role      string `json:"role"`
// 		IsActive  *bool  `json:"is_active,omitempty"`
// 	}

// 	if err := web.UnmarshalJSON(r, &req); err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}

// 	ctx := r.Context()
// 	if err := h.Service.UpdateCustomer(ctx, id, req.FirstName, req.LastName, req.Email, req.Role, req.IsActive); err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, map[string]string{
// 		"message": "Customer updated successfully",
// 	})
// }

// func (h *CustomerHandler) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.adminOnly(w, r) {
// 		return
// 	}

// 	idStr := mux.Vars(r)["id"]
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer ID")
// 		return
// 	}

// 	ctx := r.Context()
// 	if err := h.Service.DeleteCustomer(ctx, id); err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, map[string]string{
// 		"message": "Customer deleted successfully",
// 	})
// }

// func RegisterCustomerRoutes(router *mux.Router, customerHandler *CustomerHandler) {
// 	router.HandleFunc("/customers", customerHandler.GetCustomersHandler).Methods("GET")
// 	router.HandleFunc("/customers/{id}", customerHandler.GetCustomerByIDHandler).Methods("GET")
// 	router.HandleFunc("/customers", customerHandler.CreateCustomerHandler).Methods("POST")
// 	router.HandleFunc("/customers/{id}", customerHandler.UpdateCustomerHandler).Methods("PUT")
// 	router.HandleFunc("/customers/{id}", customerHandler.DeleteCustomerHandler).Methods("DELETE")
// }

package handler

import (
	"banking-app/component/customer/service"
	"banking-app/middleware"
	"banking-app/utils"
	"banking-app/web"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	Service *service.CustomerService
}

func NewCustomerHandler(cs *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{Service: cs}
}

// adminOnly helper
func (h *CustomerHandler) adminOnly(w http.ResponseWriter, r *http.Request) bool {
	claims, ok := middleware.GetUserClaims(r)
	if !ok || claims.Role != "admin" {
		web.RespondErrorMessage(w, http.StatusForbidden, "only admin can perform this action")
		return false
	}
	return true
}

// ---------------- CREATE CUSTOMER ----------------
func (h *CustomerHandler) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if !h.adminOnly(w, r) {
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Role      string `json:"role"`
		IsActive  *bool  `json:"is_active,omitempty"`
	}

	if err := web.UnmarshalJSON(r, &req); err != nil {
		web.RespondError(w, err)
		return
	}

	if req.Role == "" {
		req.Role = "customer"
	}

	ctx := r.Context()
	err := h.Service.CreateCustomer(ctx, req.FirstName, req.LastName, req.Email, req.Password, req.Role)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "Customer created successfully",
	})
}

func (h *CustomerHandler) GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filters := map[string]string{
		"first_name": query.Get("first_name"),
		"last_name":  query.Get("last_name"),
		"email":      query.Get("email"),
		"role":       query.Get("role"),
	}

	pagination := utils.GetPaginationParams(r, 10, 1)
	offset := (pagination.Offset - 1) * pagination.Limit

	// total count
	total, err := h.Service.CountCustomers(r.Context(), filters)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fetch paginated customers
	customers, err := h.Service.ListCustomers(r.Context(), pagination.Limit, offset, filters)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, utils.PaginatedResponse(customers, total, pagination.Limit, offset))
}

// ---------------- GET CUSTOMER BY ID ----------------
func (h *CustomerHandler) GetCustomerByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	ctx := r.Context()
	customer, err := h.Service.GetCustomerByID(ctx, id)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	// Respond with the full customer object (exclude password)
	response := map[string]interface{}{
		"customer_id": customer.CustomerID,
		"first_name":  customer.FirstName,
		"last_name":   customer.LastName,
		"email":       customer.Email,
		"role":        customer.Role,
		"is_active":   customer.IsActive,
	}
	web.RespondJSON(w, http.StatusOK, response)
}

// ---------------- UPDATE CUSTOMER ----------------
func (h *CustomerHandler) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if !h.adminOnly(w, r) {
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		IsActive  *bool  `json:"is_active,omitempty"`
	}

	if err := web.UnmarshalJSON(r, &req); err != nil {
		web.RespondError(w, err)
		return
	}

	ctx := r.Context()
	if err := h.Service.UpdateCustomer(ctx, id, req.FirstName, req.LastName, req.Email, req.Role, req.IsActive); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Customer updated successfully",
	})
}

// ---------------- DELETE CUSTOMER ----------------
func (h *CustomerHandler) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if !h.adminOnly(w, r) {
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	ctx := r.Context()
	if err := h.Service.DeleteCustomer(ctx, id); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Customer deleted successfully",
	})
}
func RegisterCustomerRoutes(router *mux.Router, customerHandler *CustomerHandler) {
	// All routes admin-only
	router.Handle("/customers", middleware.AdminOnly(http.HandlerFunc(customerHandler.GetCustomersHandler))).Methods("GET")
	router.Handle("/customers/{id}", middleware.AdminOnly(http.HandlerFunc(customerHandler.GetCustomerByIDHandler))).Methods("GET")
	router.Handle("/customers", middleware.AdminOnly(http.HandlerFunc(customerHandler.CreateCustomerHandler))).Methods("POST")
	router.Handle("/customers/{id}", middleware.AdminOnly(http.HandlerFunc(customerHandler.UpdateCustomerHandler))).Methods("PUT")
	router.Handle("/customers/{id}", middleware.AdminOnly(http.HandlerFunc(customerHandler.DeleteCustomerHandler))).Methods("DELETE")
}
