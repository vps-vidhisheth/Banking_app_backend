package handler

import (
	"net/http"

	"banking-app/component/auth/service"
	"banking-app/web"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := web.UnmarshalJSON(r, &req); err != nil {
		web.RespondError(w, err)
		return
	}

	if err := h.AuthService.Authenticate(r.Context(), req.Email, req.Password); err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, err.Error())
		return
	}

	customer, err := h.AuthService.CustomerService.GetCustomerByEmail(r.Context(), req.Email)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.AuthService.GenerateToken(customer)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, loginResponse{Token: token})
}

func RegisterAuthRoutes(router *http.ServeMux, h *AuthHandler) {
	router.HandleFunc("/auth/login", h.LoginHandler)
}
