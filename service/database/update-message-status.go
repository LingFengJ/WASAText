package database

import (
	"time"
)

func (db *appdbimpl) UpdateMessageStatus(messageID, userID, status string) error {
	// Validate status
	if status != "received" && status != "read" {
		return ErrInvalidMessageStatus
	}

	// Check if message exists
	var exists bool
	err := db.c.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM messages WHERE id = ?)",
		messageID,
	).Scan(&exists)
	if err != nil {
		return ErrDatabaseError
	}
	if !exists {
		return ErrMessageNotFound
	}

	now := time.Now()

	// Update or insert status
	_, err = db.c.Exec(`
        INSERT INTO message_status (message_id, user_id, status, updated_at)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(message_id, user_id) 
        DO UPDATE SET status = ?, updated_at = ?`,
		messageID, userID, status, now,
		status, now,
	)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}
