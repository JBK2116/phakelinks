package link

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

// LinkConn holds the database connection for link-related queries.
type LinkConn struct {
	logger *slog.Logger
	db     *sql.DB
}

// NewLinkConn() creates a new LinkConn with the provided database connection.
func NewLinkConn(db *sql.DB, logger *slog.Logger) *LinkConn {
	return &LinkConn{
		db:     db,
		logger: logger,
	}
}

// RegisterRoutes() registers all routes for the LinkConn struct
func (linkConn *LinkConn) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/links", linkConn.handleCreateLink).Methods("POST")
}

// handleCreateLink() handles the business logic for creating a new link
func (linkConn *LinkConn) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Header", r.Header)
	fmt.Println("Method", r.Method)
	fmt.Println("Body: ", r.Body)
}
