package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

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
	errCh := make(chan error, 2)

	mainServer := NewAPIServer(fmt.Sprintf(":%s", configs.Envs.PublicPort), logger, db)
	redirectServer := NewAPIServer(fmt.Sprintf(":%s", configs.Envs.RedirectPort), logger, db)
	logger.Info("Main Server running", slog.String("host", configs.Envs.PublicHost), slog.String("port", configs.Envs.PublicPort))
	logger.Info("Redirect Server running", slog.String("host", configs.Envs.RedirectHost), slog.String("port", configs.Envs.RedirectPort))
	go func() { errCh <- mainServer.Run() }()
	go func() { errCh <- redirectServer.RunRedirect() }()
	if err := <-errCh; err != nil && err != http.ErrServerClosed {
		logger.Error("Server Failed. Shutting Down ...", slog.String("error", err.Error()))
		close(errCh)
		os.Exit(1)
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
	linkConn := link.NewLinkConn(server.logger, server.db)
	linkConn.RegisterRoutes(subrouter)
	if !configs.Envs.IsDev {
		fs := http.FileServer(http.Dir("/home/jovbk/phakelinks/frontend/dist"))
		router.PathPrefix("/").Handler(fs)
	}
	return http.ListenAndServe(server.address, wrappedRouter)
}

func (server *APIServer) RunRedirect() error {
	router := mux.NewRouter()
	wrappedRouter := middleware.StripTrailingSlashMiddleware(router)
	linkConn := link.NewLinkConn(server.logger, server.db)
	linkConn.RegisterRedirectRoutes(router)
	return http.ListenAndServe(server.address, wrappedRouter)
}
