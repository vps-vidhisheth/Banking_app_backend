package handler

import (
	"banking-app/middleware"
	"banking-app/service"
	"banking-app/utils"
	"banking-app/web"
	"math"
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

// Create a new customer (Admin only)
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

	cust, err := h.Service.CreateCustomer(req.FirstName, req.LastName, req.Email, req.Password, req.Role)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.IsActive != nil {
		_ = h.Service.UpdateCustomer(cust.CustomerID, "", "", "", "", req.IsActive)
	}

	web.RespondJSON(w, http.StatusCreated, cust)
}

// Get all customers with pagination + filtering
func (h *CustomerHandler) GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	lastName := r.URL.Query().Get("last_name")
	pagination := utils.GetPaginationParams(r, 10, 1)

	customers, total, err := h.Service.GetAllCustomersPaginated(lastName, pagination.Offset, pagination.Limit)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]interface{}{
		"data":       customers,
		"total":      total,
		"page":       pagination.Offset,
		"page_size":  pagination.Limit,
		"total_page": int(math.Ceil(float64(total) / float64(pagination.Limit))),
	}

	web.RespondJSON(w, http.StatusOK, resp)
}

// Get customer by UUID
func (h *CustomerHandler) GetCustomerByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	cust, err := h.Service.GetCustomerByID(id)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, cust)
}

// Update customer (Admin only)
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

	if err := h.Service.UpdateCustomer(id, req.FirstName, req.LastName, req.Email, req.Role, req.IsActive); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Customer updated successfully",
	})
}

// Delete customer (Admin only)
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

	if err := h.Service.DeleteCustomer(id); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Customer deleted successfully",
	})
}

// RegisterCustomerRoutes adds all customer routes to the router
func RegisterCustomerRoutes(router *mux.Router, customerHandler *CustomerHandler) {
	router.HandleFunc("/customers", customerHandler.GetCustomersHandler).Methods("GET")
	router.HandleFunc("/customers/{id}", customerHandler.GetCustomerByIDHandler).Methods("GET")
	router.HandleFunc("/customers", customerHandler.CreateCustomerHandler).Methods("POST")
	router.HandleFunc("/customers/{id}", customerHandler.UpdateCustomerHandler).Methods("PUT")
	router.HandleFunc("/customers/{id}", customerHandler.DeleteCustomerHandler).Methods("DELETE")
}
