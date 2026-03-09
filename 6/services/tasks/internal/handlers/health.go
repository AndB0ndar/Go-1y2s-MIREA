package handlers

import (
    "net/http"
    "github.com/sirupsen/logrus"
)

func HealthCheck(log *logrus.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }
}
