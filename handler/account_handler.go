package handler

import (
	"banking-app/middleware"
	"banking-app/service"
	"banking-app/utils"
	"banking-app/web"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func RegisterAccountRoutes(router *mux.Router, h *AccountHandler) {
	router.HandleFunc("/accounts", h.CreateAccountHandler).Methods("POST")
	router.HandleFunc("/accounts", h.ListAccountsHandler).Methods("GET")
	router.HandleFunc("/accounts/{id}", h.GetAccountHandler).Methods("GET")
	router.HandleFunc("/accounts/{id}", h.UpdateAccountHandler).Methods("PUT")
	router.HandleFunc("/accounts/{id}", h.DeleteAccountHandler).Methods("DELETE")
	router.HandleFunc("/accounts/{id}/deposit", h.DepositHandler).Methods("POST")
	router.HandleFunc("/accounts/{id}/withdraw", h.WithdrawHandler).Methods("POST")
	router.HandleFunc("/accounts/transfer", h.TransferHandler).Methods("POST")
}

// Staff-only middleware
func (h *AccountHandler) staffOnly(w http.ResponseWriter, r *http.Request) bool {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
		return false
	}
	if claims.Role != "staff" {
		web.RespondErrorMessage(w, http.StatusForbidden, "only staff can access this resource")
		return false
	}
	return true
}
func (h *AccountHandler) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	pagination := utils.GetPaginationParams(r, 10, 0)

	var customerID, bankID uuid.UUID
	if v := r.URL.Query().Get("customer_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			customerID = id
		}
	}
	if v := r.URL.Query().Get("bank_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			bankID = id
		}
	}

	total, err := h.service.CountAccounts(customerID, bankID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	accounts, err := h.service.ListAccounts(pagination.Limit, pagination.Offset, customerID, bankID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, utils.PaginatedResponse(accounts, total, pagination.Limit, pagination.Offset))
}

func (h *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	var payload struct {
		CustomerID string `json:"customer_id"`
		BankID     string `json:"bank_id"`
	}
	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
		return
	}

	customerID, err := uuid.Parse(payload.CustomerID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer_id")
		return
	}

	bankID, err := uuid.Parse(payload.BankID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank_id")
		return
	}

	acc, err := h.service.CreateAccount(customerID, bankID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusCreated, acc)
}

func (h *AccountHandler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	acc, err := h.service.GetAccountByID(id)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, acc)
}

func (h *AccountHandler) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var payload struct {
		BankID string `json:"bank_id"`
	}
	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
		return
	}

	bankID, err := uuid.Parse(payload.BankID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank_id")
		return
	}

	acc, err := h.service.GetAccountByID(id)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	acc.BankID = bankID
	if err := h.service.UpdateAccount(acc); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, acc)
}

func (h *AccountHandler) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	if err := h.service.DeleteAccount(id); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "account deleted successfully"})
}

func (h *AccountHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var payload struct {
		CustomerID string  `json:"customer_id"`
		Amount     float64 `json:"amount"`
	}
	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
		return
	}

	customerID, err := uuid.Parse(payload.CustomerID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer_id")
		return
	}

	if err := h.service.Deposit(accountID, customerID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "deposit successful"})
}

func (h *AccountHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var payload struct {
		CustomerID string  `json:"customer_id"`
		Amount     float64 `json:"amount"`
	}
	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
		return
	}

	customerID, err := uuid.Parse(payload.CustomerID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer_id")
		return
	}

	if err := h.service.Withdraw(accountID, customerID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "withdraw successful"})
}

func (h *AccountHandler) TransferHandler(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	var payload struct {
		FromAccountID  string  `json:"from_account_id"`
		ToAccountID    string  `json:"to_account_id"`
		FromCustomerID string  `json:"from_customer_id"`
		ToCustomerID   string  `json:"to_customer_id"`
		Amount         float64 `json:"amount"`
	}
	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
		return
	}

	fromAccID, err := uuid.Parse(payload.FromAccountID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid from_account_id")
		return
	}

	toAccID, err := uuid.Parse(payload.ToAccountID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid to_account_id")
		return
	}

	fromCustID, err := uuid.Parse(payload.FromCustomerID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid from_customer_id")
		return
	}

	toCustID, err := uuid.Parse(payload.ToCustomerID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid to_customer_id")
		return
	}

	if err := h.service.Transfer(fromAccID, toAccID, fromCustID, toCustID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "transfer successful"})
}
