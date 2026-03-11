package middleware

import (
	"net/http"
)

func InstanceID(instanceID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Instance-ID", instanceID)
			next.ServeHTTP(w, r)
		})
	}
}
