package httpx

import (
	"net/http"
)

func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func BadRequest(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusUnauthorized)
}

func InternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
}
