package utils

import (
	"net/http"

	jsonUtils "github.com/rubeniriso/portfolio-backend/utils/json"
)

type ApiError struct {
	Error string `json:"error"`
}

func permissionDenied(w http.ResponseWriter) {
	jsonUtils.WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}
func badRequest(w http.ResponseWriter) {
	jsonUtils.WriteJSON(w, http.StatusForbidden, ApiError{Error: "bad request"})
}
