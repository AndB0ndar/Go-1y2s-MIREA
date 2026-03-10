package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"app/services/auth/internal/service"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func Login(
	authService *service.AuthService, log *logrus.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.WithError(err).Warn("Invalid login request")
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		token, err := authService.Login(req.Username, req.Password)
		if err != nil {
			log.WithError(err).Warn("Login failed")
			http.Error(
				w, `{"error":"invalid credentials"}`, http.StatusUnauthorized,
			)
			return
		}

		csrfToken := uuid.New().String()

		// Setting session cookie (HttpOnly, Secure, SameSite=Lax)
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true, // necessarily for HTTPS
			SameSite: http.SameSiteLaxMode,
			MaxAge:   3600, // 1 h
		})

		// Setting csrf cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Path:     "/",
			HttpOnly: false,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   3600,
		})

		resp := LoginResponse{
			AccessToken: token,
			TokenType:   "Bearer",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
