package middleware

import (
	"net/http"
)

// AdminOnly middleware allows only admins
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserClaims(r)
		if !ok {
			http.Error(w, "user not found in context", http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" {
			http.Error(w, "forbidden: only admins allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// StaffOnly middleware allows only staff
func StaffOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserClaims(r)
		if !ok {
			http.Error(w, "user not found in context", http.StatusUnauthorized)
			return
		}

		if claims.Role != "staff" {
			http.Error(w, "forbidden: only staff allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
