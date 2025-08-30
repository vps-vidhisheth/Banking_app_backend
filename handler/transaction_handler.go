package handler

import (
	"banking-app/model"
	"banking-app/service"
	"banking-app/utils"
	"banking-app/web"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid request body")
		return
	}

	accountID, err := uuid.Parse(payload.AccountID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account_id format")
		return
	}

	if err := h.service.RecordDeposit(accountID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "Deposit successful"})
}

func (h *TransactionHandler) WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid request body")
		return
	}

	accountID, err := uuid.Parse(payload.AccountID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account_id format")
		return
	}

	if err := h.service.RecordWithdrawal(accountID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "Withdrawal successful"})
}

func (h *TransactionHandler) TransferHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FromAccount string  `json:"from_account"`
		ToAccount   string  `json:"to_account"`
		Amount      float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fromID, err := uuid.Parse(payload.FromAccount)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid from_account format")
		return
	}

	toID, err := uuid.Parse(payload.ToAccount)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid to_account format")
		return
	}

	if err := h.service.RecordTransfer(fromID, toID, payload.Amount); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "Transfer successful"})
}

func (h *TransactionHandler) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	accountIDStr := mux.Vars(r)["id"]
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account ID format")
		return
	}

	transactions, err := h.service.GetTransactionsByAccount(accountID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) GetNetTransfersHandler(w http.ResponseWriter, r *http.Request) {
	accountIDStr := mux.Vars(r)["id"]
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account ID format")
		return
	}

	net, err := h.service.GetNetTransfers(accountID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]float64{"net_transfers": net})
}

func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	pagination := utils.GetPaginationParams(r, 2, 0)

	var transactions []model.Transaction
	var err error

	accountParam := r.URL.Query().Get("account_id")
	if accountParam != "" {
		accountID, parseErr := uuid.Parse(accountParam)
		if parseErr != nil {
			web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account_id format")
			return
		}
		transactions, err = h.service.GetTransactionsByAccount(accountID)
	} else {
		transactions, err = h.service.GetAllTransactions()
	}

	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, utils.PaginatedResponse(transactions, int64(len(transactions)), pagination.Limit, pagination.Offset))
}

func RegisterTransactionRoutes(r *mux.Router, h *TransactionHandler) {
	r.HandleFunc("/transactions/account/{id}", h.GetTransactionsHandler).Methods("GET")
	r.HandleFunc("/transactions/account/{id}/net", h.GetNetTransfersHandler).Methods("GET")
	r.HandleFunc("/transactions", h.GetAllTransactions).Methods("GET")
	r.HandleFunc("/transactions/deposit", h.DepositHandler).Methods("POST")
	r.HandleFunc("/transactions/withdrawal", h.WithdrawalHandler).Methods("POST")
	r.HandleFunc("/transactions/transfer", h.TransferHandler).Methods("POST")
}
