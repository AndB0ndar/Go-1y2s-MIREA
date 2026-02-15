package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type verifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Verify(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(verifyResponse{Valid: false, Error: "missing token"})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(verifyResponse{Valid: false, Error: "invalid authorization format"})
		return
	}
	token := parts[1]
	if token != "demo-token" { // упрощённая проверка
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(verifyResponse{Valid: false, Error: "invalid token"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(verifyResponse{Valid: true, Subject: "student"})
}
