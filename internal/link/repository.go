package link

import (
	"database/sql"
)

// InsertLink() Inserts a link into the database
func InsertLink(db *sql.DB, link string, fakeLink string) error {
	insertStmt := `INSERT INTO links (link, fakelink) VALUES ($1, $2)`
	_, err := db.Exec(insertStmt, link, fakeLink)
	return err
}

// getLink() Retreives a link from the database with a matching target string
func GetLink(db *sql.DB, target string) (string, error) {
	getStmt := `SELECT link FROM links WHERE fakelink = ($1) LIMIT 1`
	var link string
	err := db.QueryRow(getStmt, target).Scan(&link)
	if err != nil {
		return "", err
	}
	return link, nil
}
