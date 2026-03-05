package handlers

import (
	"encoding/json"
	"net/http"

	"app/services/auth/internal/service"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	svc := service.NewAuthService()
	token, err := svc.Login(req.Username, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		} else {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		}
		return
	}
	resp := loginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
