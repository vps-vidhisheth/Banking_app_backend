// package handler

// import (
// 	"banking-app/middleware"
// 	"banking-app/service"
// 	"banking-app/utils"
// 	"banking-app/web"
// 	"net/http"

// 	"github.com/google/uuid"
// 	"github.com/gorilla/mux"
// )

// type LedgerHandler struct {
// 	LedgerService *service.LedgerService
// }

// func NewLedgerHandler(ledgerService *service.LedgerService) *LedgerHandler {
// 	return &LedgerHandler{LedgerService: ledgerService}
// }

// // staffOnly helper
// func (h *LedgerHandler) staffOnly(w http.ResponseWriter, r *http.Request) bool {
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

// func (h *LedgerHandler) GetAllLedgers(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}

// 	pagination := utils.GetPaginationParams(r, 2, 0)

// 	var accountID uuid.UUID
// 	if v := r.URL.Query().Get("account_id"); v != "" {
// 		id, err := uuid.Parse(v)
// 		if err != nil {
// 			web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account_id format")
// 			return
// 		}
// 		accountID = id
// 	}

// 	ledgers, total, err := h.LedgerService.GetAllLedgers(accountID, pagination.Limit, pagination.Offset)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, utils.PaginatedResponse(ledgers, total, pagination.Limit, pagination.Offset))
// }

// func (h *LedgerHandler) GetLedger(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}

// 	id := web.ParseUUIDParam(r, "id")
// 	if id == uuid.Nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid ledger id")
// 		return
// 	}

// 	ledger, err := h.LedgerService.GetLedger(id)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, ledger)
// }

// func (h *LedgerHandler) GetNetBankTransfer(w http.ResponseWriter, r *http.Request) {
// 	if !h.staffOnly(w, r) {
// 		return
// 	}

// 	bankFromStr := r.URL.Query().Get("bank_from_id")
// 	bankToStr := r.URL.Query().Get("bank_to_id")

// 	bankFromID, err := uuid.Parse(bankFromStr)
// 	if err != nil || bankFromID == uuid.Nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank_from_id")
// 		return
// 	}

// 	bankToID, err := uuid.Parse(bankToStr)
// 	if err != nil || bankToID == uuid.Nil {
// 		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank_to_id")
// 		return
// 	}

// 	net, err := h.LedgerService.GetNetBankTransfer(bankFromID, bankToID)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, map[string]interface{}{
// 		"bank_from_id": bankFromID,
// 		"bank_to_id":   bankToID,
// 		"net_transfer": net,
// 	})
// }

// func RegisterLedgerRoutes(r *mux.Router, h *LedgerHandler) {
// 	r.HandleFunc("/ledgers", h.GetAllLedgers).Methods("GET")
// 	r.HandleFunc("/ledgers/{id}", h.GetLedger).Methods("GET")
// 	r.HandleFunc("/ledgers/net-transfer", h.GetNetBankTransfer).Methods("GET")
// }

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

type LedgerHandler struct {
	LedgerService *service.LedgerService
}

func NewLedgerHandler(ledgerService *service.LedgerService) *LedgerHandler {
	return &LedgerHandler{LedgerService: ledgerService}
}

// staffOnly helper
func (h *LedgerHandler) staffOnly(w http.ResponseWriter, r *http.Request) bool {
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

// GET /ledgers?account_id=<optional>
func (h *LedgerHandler) GetAllLedgers(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	pagination := utils.GetPaginationParams(r, 2, 0)

	var accountID *uuid.UUID
	if v := r.URL.Query().Get("account_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			web.RespondErrorMessage(w, http.StatusBadRequest, "invalid account_id format")
			return
		}
		accountID = &id
	}

	ledgers, total, err := h.LedgerService.GetAllLedgers(accountID, pagination.Limit, pagination.Offset)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, utils.PaginatedResponse(ledgers, total, pagination.Limit, pagination.Offset))
}

// GET /ledgers/{id}
func (h *LedgerHandler) GetLedger(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	id := web.ParseUUIDParam(r, "id")
	if id == uuid.Nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid ledger id")
		return
	}

	ledger, err := h.LedgerService.GetLedger(id)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, ledger)
}

// GET /ledgers/net-transfer?bank_from_id=...&bank_to_id=...
func (h *LedgerHandler) GetNetBankTransfer(w http.ResponseWriter, r *http.Request) {
	if !h.staffOnly(w, r) {
		return
	}

	bankFromStr := r.URL.Query().Get("bank_from_id")
	bankToStr := r.URL.Query().Get("bank_to_id")

	bankFromID, err := uuid.Parse(bankFromStr)
	if err != nil || bankFromID == uuid.Nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank_from_id")
		return
	}

	bankToID, err := uuid.Parse(bankToStr)
	if err != nil || bankToID == uuid.Nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank_to_id")
		return
	}

	net, err := h.LedgerService.GetNetBankTransfer(bankFromID, bankToID)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"bank_from_id": bankFromID,
		"bank_to_id":   bankToID,
		"net_transfer": net,
	})
}
func RegisterLedgerRoutes(r *mux.Router, h *LedgerHandler) {
	r.HandleFunc("/ledgers/net-transfer", h.GetNetBankTransfer).Methods("GET") // put first
	r.HandleFunc("/ledgers", h.GetAllLedgers).Methods("GET")
	r.HandleFunc("/ledgers/{id}", h.GetLedger).Methods("GET")
}
