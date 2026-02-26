package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type ctxKeyRequestID struct{}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	if val, ok := ctx.Value(ctxKeyRequestID{}).(string); ok {
		return val
	}
	return ""
}
