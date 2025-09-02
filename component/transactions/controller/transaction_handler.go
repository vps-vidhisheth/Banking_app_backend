// package handler

// import (
// 	"banking-app/component/transactions/service"
// 	"banking-app/web"
// 	"encoding/json"
// 	"net/http"

// 	"github.com/google/uuid"
// 	"github.com/gorilla/mux"
// )

// type TransactionHandler struct {
// 	service *service.TransactionService
// }

// func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
// 	return &TransactionHandler{service: service}
// }

// // DepositHandler
// func (h *TransactionHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload struct {
// 		AccountID string  `json:"account_id"`
// 		Amount    float64 `json:"amount"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid request body")
// 		return
// 	}

// 	accountID, err := uuid.Parse(payload.AccountID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account_id format")
// 		return
// 	}

// 	if err := h.service.RecordDeposit(r.Context(), accountID, payload.Amount, nil); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "Deposit successful"})
// }

// // WithdrawalHandler
// func (h *TransactionHandler) WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload struct {
// 		AccountID string  `json:"account_id"`
// 		Amount    float64 `json:"amount"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid request body")
// 		return
// 	}

// 	accountID, err := uuid.Parse(payload.AccountID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account_id format")
// 		return
// 	}

// 	if err := h.service.RecordWithdrawal(r.Context(), accountID, payload.Amount, nil); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "Withdrawal successful"})
// }

// // TransferHandler
// func (h *TransactionHandler) TransferHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload struct {
// 		FromAccount string  `json:"from_account"`
// 		ToAccount   string  `json:"to_account"`
// 		Amount      float64 `json:"amount"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid request body")
// 		return
// 	}

// 	fromID, err := uuid.Parse(payload.FromAccount)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid from_account format")
// 		return
// 	}

// 	toID, err := uuid.Parse(payload.ToAccount)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid to_account format")
// 		return
// 	}

// 	if err := h.service.RecordTransfer(r.Context(), fromID, toID, payload.Amount, nil); err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "Transfer successful"})
// }

// // GetTransactionsHandler (single account)
// func (h *TransactionHandler) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
// 	page := 1
// 	limit := 10

// 	accountIDStr := mux.Vars(r)["id"]
// 	accountID, err := uuid.Parse(accountIDStr)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account ID format")
// 		return
// 	}

// 	if r.URL.Query().Get("page") != "" {
// 		// parse page if provided
// 	}

// 	if r.URL.Query().Get("limit") != "" {
// 		// parse limit if provided
// 	}

// 	resp, err := h.service.GetTransactionsByAccount(r.Context(), accountID, page, limit)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, resp)
// }

// // GetAllTransactions (optional account filter)
// func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
// 	page := 1
// 	limit := 10

// 	if r.URL.Query().Get("page") != "" {
// 		// parse page if provided
// 	}

// 	if r.URL.Query().Get("limit") != "" {
// 		// parse limit if provided
// 	}

// 	accountParam := r.URL.Query().Get("account_id")

// 	var resp map[string]interface{}
// 	var err error

// 	if accountParam != "" {
// 		accountID, parseErr := uuid.Parse(accountParam)
// 		if parseErr != nil {
// 			web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account_id format")
// 			return
// 		}
// 		resp, err = h.service.GetTransactionsByAccount(r.Context(), accountID, page, limit)
// 	} else {
// 		resp, err = h.service.GetAllTransactions(r.Context(), page, limit)
// 	}

// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, resp)
// }

// // RegisterTransactionRoutes registers all transaction routes
//
//	func RegisterTransactionRoutes(router *mux.Router, h *TransactionHandler) {
//		router.HandleFunc("/transactions/deposit", h.DepositHandler).Methods("POST")
//		router.HandleFunc("/transactions/withdraw", h.WithdrawalHandler).Methods("POST")
//		router.HandleFunc("/transactions/transfer", h.TransferHandler).Methods("POST")
//		router.HandleFunc("/transactions/{id}", h.GetTransactionsHandler).Methods("GET")
//		router.HandleFunc("/transactions", h.GetAllTransactions).Methods("GET")
//	}

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
