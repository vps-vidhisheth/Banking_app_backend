package handler

import (
	"banking-app/component/transactions/service"
	"banking-app/utils"
	"banking-app/web"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// GetTransactionsHandler - single transaction by ID
func (h *TransactionHandler) GetTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid transaction ID format")
		return
	}

	err = h.service.GetTransactionByID(r.Context(), id)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "transaction exists"})
}

// GetTransactionsHandler - list transactions with optional filters
func (h *TransactionHandler) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	pagination := utils.GetPaginationParams(r, 10, 0)
	limit := pagination.Limit
	offset := pagination.Offset

	var accountID *uuid.UUID
	if v := r.URL.Query().Get("account_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account_id")
			return
		}
		accountID = &id
	}

	txType := r.URL.Query().Get("type")
	note := r.URL.Query().Get("note")

	var startDate, endDate *time.Time
	if v := r.URL.Query().Get("start_date"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			web.RespondErrorMessage(w, http.StatusBadRequest, "invalid start_date format, must be RFC3339")
			return
		}
		startDate = &t
	}
	if v := r.URL.Query().Get("end_date"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			web.RespondErrorMessage(w, http.StatusBadRequest, "invalid end_date format, must be RFC3339")
			return
		}
		endDate = &t
	}

	results, total, err := h.service.GetTransactions(r.Context(), accountID, txType, note, startDate, endDate, limit, offset)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, utils.PaginatedResponse(results, total, limit, offset))
}

// RegisterTransactionRoutes
func RegisterTransactionRoutes(router *mux.Router, h *TransactionHandler) {
	router.HandleFunc("/transactions/{id}", h.GetTransactionByIDHandler).Methods("GET")
	router.HandleFunc("/transactions", h.GetTransactionsHandler).Methods("GET")
}
