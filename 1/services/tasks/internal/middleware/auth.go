package middleware

import (
	"app/services/tasks/internal/auth"
	"app/shared/middleware"
	"net/http"
	"strings"
)

func Auth(authClient *auth.Client) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}
			token := parts[1]
			reqID := middleware.GetRequestID(r.Context())

			err := authClient.VerifyToken(r.Context(), token, reqID)
			if err != nil {
				if err.Error() == "unauthorized" {
					http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error":"authorization service unavailable"}`, http.StatusServiceUnavailable)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}
