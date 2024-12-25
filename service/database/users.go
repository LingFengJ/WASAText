package database

import (
    "github.com/gofrs/uuid"
    "time"
    "fmt"
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
        return "", fmt.Errorf("error generating identifier: %w", err)
    }

    // Generate unique ID
    id, err := uuid.NewV4()
    if err != nil {
        return "", fmt.Errorf("error generating id: %w", err)
    }

    now := time.Now()

    _, err = db.c.Exec(`
        INSERT INTO users (
            id, username, password, identifier, created_at, modified_at
        ) VALUES (?, ?, ?, ?, ?, ?)`,
        id.String(), username, password, identifier.String(), now, now)
    
    if err != nil {
        // Log the specific error for debugging
        fmt.Printf("Database error: %v\n", err)
        return "", fmt.Errorf("error inserting user: %w", err)
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