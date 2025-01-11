package database

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) CheckUserExists(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *appdbimpl) CreateUser(username string, password string) (string, error) {
	// Generate unique identifier
	identifier, err := uuid.NewV4()
	if err != nil {
		var (
			ErrGeneratingIdentifier = errors.New("error generating identifier")
		)
		return "", ErrGeneratingIdentifier
	}

	// Generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		var (
			ErrGeneratingId = errors.New("error generating id")
		)
		return "", ErrGeneratingId
	}

	now := time.Now()

	_, err = db.c.Exec(`
        INSERT INTO users (
            id, username, password, identifier, created_at, modified_at
        ) VALUES (?, ?, ?, ?, ?, ?)`,
		id.String(), username, password, identifier.String(), now, now)

	if err != nil {
		var (
			ErrInsertingUser = errors.New("error inserting user")
		)
		return "", ErrInsertingUser
	}

	return identifier.String(), nil
}

func (db *appdbimpl) GetUserByCredentials(username string, password string) (string, error) {
	var identifier string
	err := db.c.QueryRow("SELECT identifier FROM users WHERE username=? AND password=?",
		username, password).Scan(&identifier)
	if err != nil {
		return "", err
	}
	return identifier, nil
}
