package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func CSRFProtection(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPatch || r.Method == http.MethodDelete {
				csrfCookie, err := r.Cookie("csrf_token")
				if err != nil {
					log.WithField("reason", "missing_csrf_cookie").Warn("CSRF check failed")
					http.Error(w, `{"error":"CSRF token missing"}`, http.StatusForbidden)
					return
				}
				headerToken := r.Header.Get("X-CSRF-Token")
				if headerToken == "" {
					log.WithField("reason", "missing_csrf_header").Warn("CSRF check failed")
					http.Error(w, `{"error":"CSRF token missing"}`, http.StatusForbidden)
					return
				}
				if csrfCookie.Value != headerToken {
					log.WithFields(logrus.Fields{
						"cookie": csrfCookie.Value,
						"header": headerToken,
					}).Warn("CSRF token mismatch")
					http.Error(w, `{"error":"CSRF token invalid"}`, http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
