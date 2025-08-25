package handler

import (
	"net/http"

	"banking-app/service"
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

	// Use web.UnmarshalJSON for request parsing
	if err := web.UnmarshalJSON(r, &req); err != nil {
		web.RespondError(w, err)
		return
	}

	// Authenticate user
	token, err := h.AuthService.Authenticate(req.Email, req.Password)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Respond with token
	web.RespondJSON(w, http.StatusOK, loginResponse{Token: token})
}
