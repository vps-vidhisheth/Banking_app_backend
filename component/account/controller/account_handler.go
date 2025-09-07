package handler

import (
	"banking-app/component/account/service"
	"banking-app/db"
	"banking-app/middleware"
	"banking-app/repository"
	"banking-app/utils"
	"banking-app/web"
	"encoding/json"
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

type TransferPayload struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description,omitempty"`
}

func RegisterAccountRoutes(router *mux.Router, h *AccountHandler) {
	router.Handle("/accounts", middleware.StaffOnly(http.HandlerFunc(h.CreateAccountHandler))).Methods("POST")
	router.Handle("/accounts", middleware.StaffOnly(http.HandlerFunc(h.ListAccountsHandler))).Methods("GET")
	router.Handle("/accounts/{id}", middleware.StaffOnly(http.HandlerFunc(h.GetAccountHandler))).Methods("GET")
	router.Handle("/accounts/{id}", middleware.StaffOnly(http.HandlerFunc(h.DeleteAccountHandler))).Methods("DELETE")
	router.Handle("/accounts/{id}/deposit", middleware.StaffOnly(http.HandlerFunc(h.DepositHandler))).Methods("POST")
	router.Handle("/accounts/{id}/withdraw", middleware.StaffOnly(http.HandlerFunc(h.WithdrawHandler))).Methods("POST")
	router.Handle("/accounts/transfer", middleware.StaffOnly(http.HandlerFunc(h.TransferHandler))).Methods("POST")
}

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

	callerUUID, err := uuid.Parse(claims.UserID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "invalid user ID in token")
		return false
	}

	if acc.CustomerID != callerUUID {
		web.RespondErrorMessage(w, http.StatusForbidden, "you can only access your own account")
		return false
	}

	return true
}

func (h *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserClaims(r)

	if !ok {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
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

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "invalid user ID")
		return
	}

	uow := repository.NewUnitOfWork(db.GetDB())
	defer uow.Rollback()

	createdAccount, err := h.service.CreateAccountWithUOW(uow, userID, bankID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := uow.Commit(); err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, "failed to commit transaction")
		return
	}

	web.RespondJSON(w, http.StatusCreated, createdAccount)
}

func (h *AccountHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

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

func (h *AccountHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

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

	if err := h.service.WithdrawWithUOW(uow, accountID, customerID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	uow.Commit()
	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "withdraw successful"})
}

func (h *AccountHandler) TransferHandler(w http.ResponseWriter, r *http.Request) {
	var payload TransferPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uow := repository.NewUnitOfWork(db.GetDB())
	defer func() {
		if r := recover(); r != nil {
			uow.Rollback()
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}()

	ctx := r.Context()

	fromAccID, err := uuid.Parse(payload.FromAccountID)
	if err != nil {
		http.Error(w, "invalid from_account_id", http.StatusBadRequest)
		return
	}
	toAccID, err := uuid.Parse(payload.ToAccountID)
	if err != nil {
		http.Error(w, "invalid to_account_id", http.StatusBadRequest)
		return
	}

	fromAcc, err := h.service.RepoGetByID(ctx, fromAccID)
	if err != nil {
		uow.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	toAcc, err := h.service.RepoGetByID(ctx, toAccID)
	if err != nil {
		uow.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.TransferWithUOW(
		uow,
		fromAcc.AccountID,
		toAcc.AccountID,
		fromAcc.CustomerID,
		toAcc.CustomerID,
		payload.Amount,
	)
	if err != nil {
		uow.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := uow.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Transfer successful"})
}

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

func (h *AccountHandler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account id")
		return
	}

	if !h.selfOnly(w, r, accountID) {
		return
	}

	acc, err := h.service.RepoGetByID(r.Context(), accountID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, acc)
}

func (h *AccountHandler) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	pagination := utils.GetPaginationParams(r, 10, 0)

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	staffID, err := uuid.Parse(claims.UserID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, "invalid user ID")
		return
	}

	searchQuery := r.URL.Query().Get("search")

	accounts, total, err := h.service.ListAccountsForStaff(r.Context(), staffID, pagination.Offset, pagination.Limit, searchQuery)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"data": accounts,
		"pagination": map[string]interface{}{
			"total":  total,
			"limit":  pagination.Limit,
			"offset": pagination.Offset,
		},
	}

	web.RespondJSON(w, http.StatusOK, response)
}
