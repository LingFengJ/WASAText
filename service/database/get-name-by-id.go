package database

import (
	"database/sql"
)

func (db *appdbimpl) GetUsernameByIdentifier(id string) (string, error) {
	var username string
	err := db.c.QueryRow(
		"SELECT username FROM users WHERE id = ?",
		id,
	).Scan(&username)

	if err == sql.ErrNoRows {
		return "", ErrUserNotFound
	}
	if err != nil {
		return "", ErrDatabaseError
	}

	return username, nil
}
