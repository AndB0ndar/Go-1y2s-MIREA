package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"app/services/tasks/internal/auth"
	"app/shared/middleware"
)

type AuthChecker interface {
	Verify(ctx context.Context, token string) (bool, error)
}

func Auth(checker AuthChecker, log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.WithField("request_id", middleware.GetRequestID(r.Context())).Warn("missing authorization header")
				http.Error(w, `{"error":"missing authorization"}`, http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.WithField("request_id", middleware.GetRequestID(r.Context())).Warn("invalid authorization format")
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}
			token := parts[1]
			valid, err := checker.Verify(r.Context(), token)
			if err != nil {
				if err == auth.ErrAuthUnavailable {
					log.WithField("request_id", middleware.GetRequestID(r.Context())).Error("auth service unavailable")
					http.Error(w, `{"error":"auth service unavailable"}`, http.StatusServiceUnavailable)
				} else {
					log.WithField("request_id", middleware.GetRequestID(r.Context())).WithError(err).Error("auth verify error")
					http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
				}
				return
			}
			if !valid {
				log.WithField("request_id", middleware.GetRequestID(r.Context())).Warn("unauthorized token")
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
