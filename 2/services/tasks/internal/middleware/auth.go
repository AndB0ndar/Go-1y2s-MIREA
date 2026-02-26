package middleware

import (
	"context"
	"net/http"
	"strings"

	"app/services/tasks/internal/auth"
)

type AuthChecker interface {
	Verify(ctx context.Context, token string) (bool, error)
}

func Auth(checker AuthChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization"}`, http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}
			token := parts[1]
			valid, err := checker.Verify(r.Context(), token)
			if err != nil {
				if err == auth.ErrAuthUnavailable {
					http.Error(w, `{"error":"auth service unavailable"}`, http.StatusServiceUnavailable)
				} else {
					http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
				}
				return
			}
			if !valid {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
