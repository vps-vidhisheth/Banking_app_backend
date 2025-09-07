package handler

import (
	"net/http"
	"strconv"
	"strings"

	"banking-app/component/banks/service"
	"banking-app/middleware"
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

	if err := h.Service.CreateBank(r.Context(), payload.Name); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusCreated, map[string]string{"message": "bank created"})
}

func (h *BankHandler) GetAllBanksHandler(w http.ResponseWriter, r *http.Request) {
	if !h.adminOnly(w, r) {
		return
	}

	limit := 10
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	name := r.URL.Query().Get("search")
	filters := map[string]string{
		"name": name,
	}

	banks, total, err := h.Service.ListBanksWithFilters(r.Context(), limit, offset, filters)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"data":   banks,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	web.RespondJSON(w, http.StatusOK, response)
}

func (h *BankHandler) GetBankByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, "invalid bank id")
		return
	}

	bank, err := h.Service.GetBankByID(r.Context(), id)
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

	if err := h.Service.UpdateBank(r.Context(), id, payload.Name); err != nil {
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

	if err := h.Service.DeleteBank(r.Context(), id); err != nil {
		web.RespondErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, map[string]string{"message": "bank deleted"})
}

func RegisterBankRoutes(router *mux.Router, h *BankHandler) {
	router.Handle("/banks", middleware.AdminOnly(http.HandlerFunc(h.CreateBankHandler))).Methods("POST")
	router.Handle("/banks/{id}", middleware.AdminOnly(http.HandlerFunc(h.GetBankByIDHandler))).Methods("GET")
	router.Handle("/banks", middleware.AdminOnly(http.HandlerFunc(h.GetAllBanksHandler))).Methods("GET")
	router.Handle("/banks/{id}", middleware.AdminOnly(http.HandlerFunc(h.UpdateBankHandler))).Methods("PUT")
	router.Handle("/banks/{id}", middleware.AdminOnly(http.HandlerFunc(h.DeleteBankHandler))).Methods("DELETE")
}
