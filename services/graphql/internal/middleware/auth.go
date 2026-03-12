package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"app/services/graphql/internal/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(authClient *auth.Client, log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}
			token := parts[1]
			valid, err := authClient.Verify(r.Context(), token)
			if err != nil {
				log.WithError(err).Error("auth verify failed")
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			if !valid {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, "user")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
