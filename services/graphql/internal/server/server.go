package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sirupsen/logrus"

	"app/services/graphql/graph/generated"
	"app/services/graphql/internal/auth"
	"app/services/graphql/internal/middleware"
	"app/services/graphql/internal/repository"
	"app/services/graphql/internal/resolvers"
	shared_middleware "app/shared/middleware"
)

func NewServer(port, authGRPCAddr string, db *sql.DB, log *logrus.Logger) *http.Server {
	repo := repository.NewPostgresRepo(db)

	// Auth
	authClient, err := auth.NewClient(authGRPCAddr, 3*time.Second)
	if err != nil {
		log.WithError(err).Warn("failed to create auth client, continuing without auth")
		authClient = nil
	}

	// GraphQL resolver
	resolver := &resolvers.Resolver{Repo: repo}
	config := generated.Config{Resolvers: resolver}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	mux := http.NewServeMux()
	//mux.Handle("/", transport.GraphiQL{Endpoint: "/query"})
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	// Middleware
	handler := shared_middleware.RequestID(mux)
	handler = shared_middleware.SecurityHeaders(handler)
	handler = shared_middleware.AccessLog(log)(handler)
	if authClient != nil {
		handler = middleware.AuthMiddleware(authClient, log)(handler)
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
