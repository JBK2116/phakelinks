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

	db, err := configs.NewPsqlConnection()
	if err != nil {
		logger.Error("Database connection error", slog.Any("error", err), slog.String("db_host", configs.Envs.DBHost), slog.Int64("db_port", configs.Envs.DBPort), slog.String("db_user", configs.Envs.DBUser), slog.String("db_name", configs.Envs.DBName), slog.String("db_password", configs.Envs.DBPassword))
		panic(err)
	}
	defer db.Close()
	logger.Info("Database successfully connected")

	server := NewAPIServer(fmt.Sprintf(":%s", configs.Envs.PublicPort), db, logger)
	logger.Info("Server running", slog.String("host", configs.Envs.PublicHost), slog.String("port", configs.Envs.PublicPort))
	if err := server.Run(); err != nil && err != http.ErrServerClosed {
		logger.Error("Error during server startup", slog.Any("error", err))
	}
}

// APIServer represents an server instance for running the application
type APIServer struct {
	address string
	db      *sql.DB
	logger  *slog.Logger
}

// NewAPIServer() returns a new APIServer instance
func NewAPIServer(address string, db *sql.DB, logger *slog.Logger) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
		logger:  logger,
	}
}

// Run() handles starting up the http server
func (server *APIServer) Run() error {
	router := mux.NewRouter()
	router.Use(middleware.StripTrailingSlashMiddleware)
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	linkConn := link.NewLinkConn(server.db, server.logger)
	linkConn.RegisterRoutes(subrouter)
	return http.ListenAndServe(server.address, router)
}
