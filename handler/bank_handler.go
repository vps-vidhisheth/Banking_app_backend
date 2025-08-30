package handler

import (
	"net/http"
	"strings"

	"banking-app/middleware"
	"banking-app/service"
	"banking-app/utils"
	"banking-app/web"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BankHandler struct {
	Service *service.BankService
}

func NewBankHandler(service *service.BankService) *BankHandler {
	return &BankHandler{Service: service}
}

func (h *BankHandler) adminOnly(w http.ResponseWriter, r *http.Request) bool {
	claims, ok := middleware.GetUserClaims(r)
	if !ok || claims.Role != "admin" {
		web.RespondErrorMessage(w, http.StatusForbidden, "only admin can perform this action")
		return false
	}
	return true
}

func (h *BankHandler) CreateBankHandler(w http.ResponseWriter, r *http.Request) {
	if !h.adminOnly(w, r) {
		return
	}

	var payload struct {
		Name string `json:"name"`
	}

	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondError(w, err)
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" {
		web.RespondErrorMessage(w, http.StatusBadRequest, "bank name cannot be empty")
		return
	}

	bank, err := h.Service.CreateBank(payload.Name)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusCreated, bank)
}

// Get all banks with pagination
func (h *BankHandler) GetAllBanksHandler(w http.ResponseWriter, r *http.Request) {

	params := utils.GetPaginationParams(r, 2, 0)

	allBanks, err := h.Service.ListBanks()
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	total := int64(len(allBanks))

	start := params.Offset
	if start > len(allBanks) {
		start = len(allBanks)
	}
	end := start + params.Limit
	if end > len(allBanks) {
		end = len(allBanks)
	}
	paginatedBanks := allBanks[start:end]

	web.RespondJSON(w, http.StatusOK, utils.PaginatedResponse(paginatedBanks, total, params.Limit, params.Offset))
}

func (h *BankHandler) GetBankByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank id")
		return
	}

	bank, err := h.Service.GetBankByID(id)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, bank)
}

func (h *BankHandler) UpdateBankHandler(w http.ResponseWriter, r *http.Request) {
	if !h.adminOnly(w, r) {
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank id")
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := web.UnmarshalJSON(r, &payload); err != nil {
		web.RespondError(w, err)
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" {
		web.RespondErrorMessage(w, http.StatusBadRequest, "bank name cannot be empty")
		return
	}

	if err := h.Service.UpdateBank(id, payload.Name); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "bank updated"})
}

func (h *BankHandler) DeleteBankHandler(w http.ResponseWriter, r *http.Request) {
	if !h.adminOnly(w, r) {
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank id")
		return
	}

	if err := h.Service.DeleteBank(id); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "bank deleted"})
}

func RegisterBankRoutes(router *mux.Router, h *BankHandler) {
	router.HandleFunc("/banks", h.CreateBankHandler).Methods("POST")
	router.HandleFunc("/banks/{id}", h.GetBankByIDHandler).Methods("GET")
	router.HandleFunc("/banks", h.GetAllBanksHandler).Methods("GET")
	router.HandleFunc("/banks/{id}", h.UpdateBankHandler).Methods("PUT")
	router.HandleFunc("/banks/{id}", h.DeleteBankHandler).Methods("DELETE")
}
