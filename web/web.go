package web

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ParseUUIDParam extracts a UUID from path variables; returns zero UUID if empty or invalid
func ParseUUIDParam(r *http.Request, key string) uuid.UUID {
	vars := mux.Vars(r)
	idStr, ok := vars[key]
	if !ok || idStr == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil
	}
	return id
}
