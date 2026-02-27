package configs

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// getNewDBConnection() returns a new database connection to the backend
func NewDBConn() (*sql.DB, error) {
	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		Envs.DBUser,
		Envs.DBPassword,
		Envs.DBHost,
		Envs.DBPort,
		Envs.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
