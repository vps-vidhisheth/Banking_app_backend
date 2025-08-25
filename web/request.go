package web

import (
	"encoding/json"
	"io"
	"net/http"

	"banking-app/apperror"
)

func UnmarshalJSON(r *http.Request, out interface{}) *apperror.AppError {
	defer r.Body.Close()

	if r.Body == nil {
		return apperror.NewValidationError("request body is empty")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return apperror.NewValidationError("failed to read request body: " + err.Error())
	}

	if len(body) == 0 {
		return apperror.NewValidationError("request body is empty")
	}

	if err := json.Unmarshal(body, out); err != nil {
		return apperror.NewValidationError("invalid JSON: " + err.Error())
	}

	return nil
}

func BindAndValidate(r *http.Request, out interface{}, validate func(interface{}) *apperror.AppError) *apperror.AppError {
	if err := UnmarshalJSON(r, out); err != nil {
		return err
	}

	if validate != nil {
		if err := validate(out); err != nil {
			return err
		}
	}

	return nil
}
