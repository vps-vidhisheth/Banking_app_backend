// package handler

// import (
// 	"net/http"

// 	"banking-app/component/auth/service"
// 	"banking-app/web"
// )

// // AuthHandler handles authentication routes
// type AuthHandler struct {
// 	AuthService *service.AuthService
// }

// // NewAuthHandler creates a new AuthHandler
// func NewAuthHandler(authService *service.AuthService) *AuthHandler {
// 	return &AuthHandler{AuthService: authService}
// }

// // loginRequest represents login payload
// type loginRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// // loginResponse represents login response
// type loginResponse struct {
// 	Token string `json:"token"`
// }

// // LoginHandler handles user login
// func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	var req loginRequest

// 	if err := web.UnmarshalJSON(r, &req); err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// pass r.Context() to AuthService
// 	token, err := h.AuthService.Authenticate(r.Context(), req.Email, req.Password)
// 	if err != nil {
// 		web.RespondErrorMessage(w, http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	web.RespondJSON(w, http.StatusOK, loginResponse{Token: token})
// }

package handler

import (
	"net/http"

	"banking-app/component/auth/service"
	"banking-app/web"
)

// AuthHandler handles authentication routes
type AuthHandler struct {
	AuthService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// loginRequest represents login payload
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginResponse represents login response
type loginResponse struct {
	Token string `json:"token"`
}

// LoginHandler handles user login
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	// Decode JSON body
	if err := web.UnmarshalJSON(r, &req); err != nil {
		web.RespondError(w, err)
		return
	}

	// Step 1: Authenticate (checks credentials, returns only error)
	if err := h.AuthService.Authenticate(r.Context(), req.Email, req.Password); err != nil {
		web.RespondErrorMessage(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Step 2: Fetch customer by email to generate token
	customer, err := h.AuthService.CustomerService.GetCustomerByEmail(r.Context(), req.Email)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Step 3: Generate JWT token
	token, err := h.AuthService.GenerateToken(customer)
	if err != nil {
		web.RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Respond with token
	web.RespondJSON(w, http.StatusOK, loginResponse{Token: token})
}

// RegisterAuthRoutes registers auth endpoints
func RegisterAuthRoutes(router *http.ServeMux, h *AuthHandler) {
	router.HandleFunc("/auth/login", h.LoginHandler)
}
