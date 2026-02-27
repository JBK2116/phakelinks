package link

import "database/sql"

func InsertLink(db *sql.DB, link string, fakeLink string) error {
	insertStmt := `INSERT INTO links (link, fakelink) VALUES ($1, $2)`
	_, err := db.Exec(insertStmt, link, fakeLink)
	return err
}
