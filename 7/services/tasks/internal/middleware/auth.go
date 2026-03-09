package middleware

import (
	"context"
	"net/http"

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
			sessionCookie, err := r.Cookie("session")
			if err != nil {
				log.WithField("reason", "missing_session_cookie").Warn("Authentication failed")
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			token := sessionCookie.Value

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
				log.WithField("token", token).Warn("Invalid session")
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
