package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/JBK2116/phakelinks/internal/configs"
	"github.com/JBK2116/phakelinks/internal/link"
	"github.com/JBK2116/phakelinks/internal/middleware"
	"github.com/gorilla/mux"
)

func main() {
	logger := configs.NewLogger(configs.Envs.IsDev)
	db, err := configs.NewDBConn()
	if err != nil {
		panic(err)
	}
	logger.Info("Database successfully connected")

	server := NewAPIServer(fmt.Sprintf(":%s", configs.Envs.PublicPort), logger, db)
	logger.Info("Server running", slog.String("host", configs.Envs.PublicHost), slog.String("port", configs.Envs.PublicPort))
	if err := server.Run(); err != nil && err != http.ErrServerClosed {
		logger.Error("Error during server startup", slog.Any("error", err))
	}
}

// APIServer represents an server instance for running the application
type APIServer struct {
	address string
	logger  *slog.Logger
	db      *sql.DB
}

// NewAPIServer() returns a new APIServer instance
func NewAPIServer(address string, logger *slog.Logger, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		logger:  logger,
		db:      db,
	}
}

// Run() handles starting up the http server
func (server *APIServer) Run() error {
	router := mux.NewRouter()
	wrappedRouter := middleware.StripTrailingSlashMiddleware(router) // router wrapping is needed here to ensure that middleware runs BEFORE matching to the path
	subrouter := router.PathPrefix("/api/v1/").Subrouter()

	linkConn := link.NewLinkConn(server.logger)
	linkConn.RegisterRoutes(subrouter)
	return http.ListenAndServe(server.address, wrappedRouter)
}
