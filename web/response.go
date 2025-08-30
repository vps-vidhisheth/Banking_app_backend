package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"banking-app/apperror"
)

func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		RespondJSON(w, appErr.Code, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
	} else {
		RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
	}
}

func RespondErrorMessage(w http.ResponseWriter, code int, msg string) {
	RespondJSON(w, code, map[string]string{"error": msg})
}

func RespondJSONWithXTotalCount(w http.ResponseWriter, code int, count int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	SetNewHeader(w, "X-Total-Count", strconv.Itoa(count))
	w.WriteHeader(code)
	w.Write(response)
}

func SetNewHeader(w http.ResponseWriter, headerName, value string) {
	w.Header().Add("Access-Control-Expose-Headers", headerName)
	w.Header().Set(headerName, value)
}
