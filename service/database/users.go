package database

import (
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
        return "", err
    }

    _, err = db.c.Exec("INSERT INTO users (username, password, identifier) VALUES (?, ?, ?)",
        username, password, identifier.String())
    if err != nil {
        return "", err
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