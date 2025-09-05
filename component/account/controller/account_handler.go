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
