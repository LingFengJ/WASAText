package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUserIDByUsername(username string) (string, error) {
	var id string
	err := db.c.QueryRow(
		"SELECT id FROM users WHERE username = ?",
		username,
	).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrUserNotFound
	}
	if err != nil {
		return "", ErrDatabaseError
	}

	return id, nil
}
