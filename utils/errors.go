package utils

import (
	"net/http"
)

type ApiError struct {
	Error string `json:"error"`
}

func PermissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}
func BadRequest(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "bad request"})
}
func DBError(w http.ResponseWriter) {
	WriteJSON(w, http.StatusConflict, ApiError{Error: "error"})
}
