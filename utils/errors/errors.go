package errors

import (
	json "kangym/utils"
	"net/http"
)

type ApiError struct {
	Error string `json:"error"`
}

func PermissionDenied(w http.ResponseWriter) {
	json.WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}
func BadRequest(w http.ResponseWriter) {
	json.WriteJSON(w, http.StatusForbidden, ApiError{Error: "bad request"})
}
func DBError(w http.ResponseWriter) {
	json.WriteJSON(w, http.StatusConflict, ApiError{Error: "error"})
}
