package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/giannuccilli/user-api/internal/domain"
)

type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

const (
	ErrCodeInvalidRequest = "INVALID_REQUEST"
	ErrCodeInvalidID      = "INVALID_ID"
	ErrCodeUserNotFound   = "USER_NOT_FOUND"
	ErrCodeEmailExists    = "EMAIL_EXISTS"
	ErrCodeInternalError  = "INTERNAL_ERROR"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, err error) {
	var status int
	var errResp ErrorResponse

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		status = http.StatusNotFound
		errResp = ErrorResponse{
			Code:    ErrCodeUserNotFound,
			Message: "User not found",
		}
	case errors.Is(err, domain.ErrEmailExists):
		status = http.StatusConflict
		errResp = ErrorResponse{
			Code:    ErrCodeEmailExists,
			Message: "Email already exists",
		}
	case errors.Is(err, domain.ErrInvalidInput):
		status = http.StatusBadRequest
		errResp = ErrorResponse{
			Code:    ErrCodeInvalidRequest,
			Message: "Invalid request data",
		}
	default:
		status = http.StatusInternalServerError
		errResp = ErrorResponse{
			Code:    ErrCodeInternalError,
			Message: "Internal server error",
		}
	}

	JSON(w, status, errResp)
}

func ErrorWithMessage(w http.ResponseWriter, status int, code, message string, details ...string) {
	errResp := ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
	JSON(w, status, errResp)
}
