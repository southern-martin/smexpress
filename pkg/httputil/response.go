package httputil

import (
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Error writes a JSON error response.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

// Created writes a 201 response with the created resource.
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, data)
}

// NoContent writes a 204 response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BadRequest writes a 400 error.
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// NotFound writes a 404 error.
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

// InternalError writes a 500 error.
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}
