package database

import "time"

func (db *appdbimpl) UpdateUsername(userID string, newUsername string) error {
    // First check if username is already taken
    var exists bool
    err := db.c.QueryRow(
        "SELECT EXISTS(SELECT 1 FROM users WHERE username = ? AND id != ?)",
        newUsername, userID,
    ).Scan(&exists)
    if err != nil {
        return ErrDatabaseError
    }
    if exists {
        return ErrUsernameTaken
    }

    // Update username
    result, err := db.c.Exec(`
        UPDATE users 
        SET username = ?, modified_at = ?
        WHERE id = ?`,
        newUsername, time.Now(), userID,
    )
    if err != nil {
        return ErrDatabaseError
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return ErrDatabaseError
    }
    if rows == 0 {
        return ErrUserNotFound
    }

    return nil
}