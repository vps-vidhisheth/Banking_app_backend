// package handler

// import (
// 	"banking-app/component/account/service"
// 	"banking-app/middleware"
// 	"banking-app/web"
// 	"net/http"

// 	"github.com/google/uuid"
// 	"github.com/gorilla/mux"
// )

// type AccountHandler struct {
// 	service *service.AccountService
// }

// func NewAccountHandler(service *service.AccountService) *AccountHandler {
// 	return &AccountHandler{service: service}
// }

// func RegisterAccountRoutes(router *mux.Router, h *AccountHandler) {
// 	router.HandleFunc("/accounts", h.CreateAccountHandler).Methods("POST")
// 	router.HandleFunc("/accounts", h.ListAccountsHandler).Methods("GET")
// 	router.HandleFunc("/accounts/{id}", h.GetAccountHandler).Methods("GET")
// 	router.HandleFunc("/accounts/{id}", h.UpdateAccountHandler).Methods("PUT")
// 	router.HandleFunc("/accounts/{id}", h.DeleteAccountHandler).Methods("DELETE")
// 	router.HandleFunc("/accounts/{id}/deposit", h.DepositHandler).Methods("POST")
// 	router.HandleFunc("/accounts/{id}/withdraw", h.WithdrawHandler).Methods("POST")
// 	router.HandleFunc("/accounts/transfer", h.TransferHandler).Methods("POST")
// }

// // Staff-only middleware
// func (h *AccountHandler) staffOnly(w http.ResponseWriter, r *http.Request) bool {
// 	claims, ok := middleware.GetUserClaims(r)
// 	if !ok {
// 		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
// 		return false
// 	}
// 	if claims.Role != "staff" {
// 		web.RespondErrorMessage(w, http.StatusForbidden, "only staff can access this resource")
// 		return false
// 	}
// 	return true
// }

// // ---------------- List Accounts ----------------
// func (h *AccountHandler) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}

// 	var customerID, bankID uuid.UUID
// 	if v := r.URL.Query().Get("customer_id"); v != "" {
// 		if id, err := uuid.Parse(v); err == nil {
// 			customerID = id
// 		}
// 	}
// 	if v := r.URL.Query().Get("bank_id"); v != "" {
// 		if id, err := uuid.Parse(v); err == nil {
// 			bankID = id
// 		}
// 	}

// 	accountsResp, err := h.service.ListAccountsWithPagination(r.Context(), r, customerID, bankID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, accountsResp)
// }

// // ---------------- Create Account ----------------
// func (h *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}
// 	ctx := r.Context()

// 	var payload struct {
// 		CustomerID string `json:"customer_id"`
// 		BankID     string `json:"bank_id"`
// 	}
// 	if err := web.UnmarshalJSON(r, &payload); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
// 		return
// 	}

// 	customerID, err := uuid.Parse(payload.CustomerID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer_id")
// 		return
// 	}

// 	bankID, err := uuid.Parse(payload.BankID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank_id")
// 		return
// 	}

// 	acc, err := h.service.CreateAccount(ctx, customerID, bankID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusCreated, acc)
// }

// // ---------------- Get Account ----------------
// func (h *AccountHandler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}
// 	ctx := r.Context()

// 	id, err := uuid.Parse(mux.Vars(r)["id"])
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
// 		return
// 	}

// 	acc, err := h.service.GetAccountByID(ctx, id)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, acc)
// }

// // ---------------- Update Account (BankID only) ----------------
// func (h *AccountHandler) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}
// 	ctx := r.Context()

// 	id, err := uuid.Parse(mux.Vars(r)["id"])
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
// 		return
// 	}

// 	var payload struct {
// 		BankID string `json:"bank_id"`
// 	}
// 	if err := web.UnmarshalJSON(r, &payload); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
// 		return
// 	}

// 	bankID, err := uuid.Parse(payload.BankID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank_id")
// 		return
// 	}

// 	acc, err := h.service.GetAccountByID(ctx, id)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
// 		return
// 	}

// 	// Update BankID
// 	acc.BankID = bankID

// 	if err := h.service.UpdateAccount(ctx, acc); err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, acc)
// }

// // ---------------- Delete Account (soft delete) ----------------
// func (h *AccountHandler) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}
// 	ctx := r.Context()

// 	id, err := uuid.Parse(mux.Vars(r)["id"])
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
// 		return
// 	}

// 	if err := h.service.SoftDeleteAccount(ctx, id); err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "account deleted successfully"})
// }

// // ---------------- Deposit (with transaction) ----------------
// func (h *AccountHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}

// 	accountID, err := uuid.Parse(mux.Vars(r)["id"])
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
// 		return
// 	}

// 	var payload struct {
// 		CustomerID string  `json:"customer_id"`
// 		Amount     float64 `json:"amount"`
// 	}
// 	if err := web.UnmarshalJSON(r, &payload); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
// 		return
// 	}

// 	customerID, err := uuid.Parse(payload.CustomerID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer_id")
// 		return
// 	}

// 	// Call service directly, no tx needed
// 	if err := h.service.Deposit(r.Context(), accountID, customerID, payload.Amount); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "deposit successful"})
// }

// // ---------------- Withdraw (with transaction) ----------------
// func (h *AccountHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}

// 	accountID, err := uuid.Parse(mux.Vars(r)["id"])
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
// 		return
// 	}

// 	var payload struct {
// 		CustomerID string  `json:"customer_id"`
// 		Amount     float64 `json:"amount"`
// 	}
// 	if err := web.UnmarshalJSON(r, &payload); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
// 		return
// 	}

// 	customerID, err := uuid.Parse(payload.CustomerID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid customer_id")
// 		return
// 	}

// 	// Call service directly
// 	if err := h.service.Withdraw(r.Context(), accountID, customerID, payload.Amount); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "withdraw successful"})
// }

// // ---------------- Transfer (with transaction) ----------------
// func (h *AccountHandler) TransferHandler(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}

// 	var payload struct {
// 		FromAccountID  string  `json:"from_account_id"`
// 		ToAccountID    string  `json:"to_account_id"`
// 		FromCustomerID string  `json:"from_customer_id"`
// 		ToCustomerID   string  `json:"to_customer_id"`
// 		Amount         float64 `json:"amount"`
// 	}
// 	if err := web.UnmarshalJSON(r, &payload); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
// 		return
// 	}

// 	fromAccID, err := uuid.Parse(payload.FromAccountID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid from_account_id")
// 		return
// 	}
// 	toAccID, err := uuid.Parse(payload.ToAccountID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid to_account_id")
// 		return
// 	}
// 	fromCustID, err := uuid.Parse(payload.FromCustomerID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid from_customer_id")
// 		return
// 	}
// 	toCustID, err := uuid.Parse(payload.ToCustomerID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid to_customer_id")
// 		return
// 	}

// 	// Call service directly
// 	if err := h.service.Transfer(r.Context(), fromAccID, toAccID, fromCustID, toCustID, payload.Amount); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

//		web.RespondJSON(w, http.StatusOK, map[string]string{"message": "transfer successful"})
//	}

package handler

import (
	"banking-app/component/account/service"
	"banking-app/db"
	"banking-app/middleware"
	"banking-app/repository"
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
	router.Handle("/accounts", middleware.StaffOnly(http.HandlerFunc(h.CreateAccountHandler))).Methods("POST")
	router.Handle("/accounts", middleware.StaffOnly(http.HandlerFunc(h.ListAccountsHandler))).Methods("GET")
	router.Handle("/accounts/{id}", middleware.StaffOnly(http.HandlerFunc(h.GetAccountHandler))).Methods("GET")
	router.Handle("/accounts/{id}", middleware.StaffOnly(http.HandlerFunc(h.UpdateAccountHandler))).Methods("PUT")
	router.Handle("/accounts/{id}", middleware.StaffOnly(http.HandlerFunc(h.DeleteAccountHandler))).Methods("DELETE")
	router.Handle("/accounts/{id}/deposit", middleware.StaffOnly(http.HandlerFunc(h.DepositHandler))).Methods("POST")
	router.Handle("/accounts/{id}/withdraw", middleware.StaffOnly(http.HandlerFunc(h.WithdrawHandler))).Methods("POST")
	router.Handle("/accounts/transfer", middleware.StaffOnly(http.HandlerFunc(h.TransferHandler))).Methods("POST")
}

// ---------------- Self Access Check ----------------
func (h *AccountHandler) selfOnly(w http.ResponseWriter, r *http.Request, accountID uuid.UUID) bool {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
		return false
	}

	acc, err := h.service.RepoGetByID(r.Context(), accountID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, "account not found")
		return false
	}

	// Parse claims.UserID string to uuid.UUID
	callerUUID, err := uuid.Parse(claims.UserID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "invalid user ID in token")
		return false
	}

	// Check if the account belongs to the logged-in user
	if acc.CustomerID != callerUUID {
		web.RespondErrorMessage(w, http.StatusForbidden, "you can only access your own account")
		return false
	}

	return true
}

// ---------------- Create Account ----------------
func (h *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Get logged-in user
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	var payload struct {
		BankID string `json:"bank_id"` // only bank_id is allowed
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

	// Use only logged-in user's ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "invalid user ID")
		return
	}

	uow := repository.NewUnitOfWork(db.GetDB())
	defer uow.Rollback()

	if err := h.service.CreateAccountWithUOW(uow, userID, bankID); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	uow.Commit()
	web.RespondJSON(w, http.StatusCreated, map[string]string{"message": "account created successfully"})
}

// ---------------- Deposit ----------------
func (h *AccountHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	// Ensure staff can only deposit into their own account
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	customerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "invalid user ID in token")
		return
	}

	if !h.selfOnly(w, r, accountID) {
		return
	}

	var payload struct {
		Amount float64 `json:"amount"`
	}
	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
		return
	}

	uow := repository.NewUnitOfWork(db.GetDB())
	defer uow.Rollback()

	if err := h.service.DepositWithUOW(uow, accountID, customerID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	uow.Commit()
	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "deposit successful"})
}

// ---------------- Withdraw ----------------
func (h *AccountHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	// Get user claims
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	customerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "invalid user ID in token")
		return
	}

	// Ensure staff can only withdraw from their own account
	if !h.selfOnly(w, r, accountID) {
		return
	}

	var payload struct {
		Amount float64 `json:"amount"`
	}
	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
		return
	}

	uow := repository.NewUnitOfWork(db.GetDB())
	defer uow.Rollback()

	if err := h.service.WithdrawWithUOW(uow, accountID, customerID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	uow.Commit()
	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "withdraw successful"})
}

func (h *AccountHandler) TransferHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FromAccountID string  `json:"from_account_id"`
		ToAccountID   string  `json:"to_account_id"`
		Amount        float64 `json:"amount"`
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

	// Get user claims
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	fromCustomerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "invalid user ID in token")
		return
	}

	// Ensure staff can only transfer from their own account
	if !h.selfOnly(w, r, fromAccID) {
		return
	}

	// Get destination account to retrieve its owner (toCustomerID)
	toAcc, err := h.service.RepoGetByID(r.Context(), toAccID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, "destination account not found")
		return
	}
	toCustomerID := toAcc.CustomerID

	uow := repository.NewUnitOfWork(db.GetDB())
	defer uow.Rollback()

	if err := h.service.TransferWithUOW(uow, fromAccID, toAccID, fromCustomerID, toCustomerID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	uow.Commit()
	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "transfer successful"})
}

// ---------------- Update Account ----------------
func (h *AccountHandler) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	if !h.selfOnly(w, r, accountID) {
		return
	}

	var payload struct {
		Balance *float64 `json:"balance"`
	}
	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Message)
		return
	}

	uow := repository.NewUnitOfWork(db.GetDB())
	defer uow.Rollback()

	acc, _ := h.service.RepoGetByID(r.Context(), accountID)
	if payload.Balance != nil {
		acc.Balance = *payload.Balance
	}

	if err := h.service.UpdateAccountWithUOW(uow, acc); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	uow.Commit()
	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "account updated successfully"})
}

// ---------------- Delete Account ----------------
func (h *AccountHandler) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	if !h.selfOnly(w, r, accountID) {
		return
	}

	uow := repository.NewUnitOfWork(db.GetDB())
	defer uow.Rollback()

	if err := h.service.SoftDeleteAccountWithUOW(uow, accountID); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	uow.Commit()
	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "account deleted successfully"})
}

// ---------------- List Accounts ----------------
func (h *AccountHandler) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "invalid user ID")
		return
	}

	// Get pagination from query params
	pagination := utils.GetPaginationParams(r, 10, 0)

	// Use service method for paginated accounts
	accounts, total, err := h.service.ListAccountsWithPaginationByUser(r.Context(), userID, pagination.Offset, pagination.Limit)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return in standard paginated format
	web.RespondJSON(w, http.StatusOK, utils.PaginatedResponse(accounts, total, pagination.Limit, pagination.Offset))
}

// ---------------- Get Account ----------------
func (h *AccountHandler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	if !h.selfOnly(w, r, accountID) {
		return
	}

	acc, _ := h.service.RepoGetByID(r.Context(), accountID)
	web.RespondJSON(w, http.StatusOK, acc)
}
