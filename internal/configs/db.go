package configs

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// NewPsqlConnection returns a new psql connection to the database
func NewPsqlConnection() (*sql.DB, error) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Envs.DBHost, Envs.DBPort, Envs.DBUser, Envs.DBPassword, Envs.DBName)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	// check db
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
