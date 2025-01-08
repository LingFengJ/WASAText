package database

import "time"

func (db *appdbimpl) AddReaction(messageID string, userID string, emoji string) error {
	// Check if reaction already exists
	var exists bool
	err := db.c.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM reactions WHERE message_id = ? AND user_id = ?)",
		messageID, userID,
	).Scan(&exists)
	if err != nil {
		return ErrDatabaseError
	}
	if exists {
		return ErrDuplicateReaction
	}

	// Add reaction
	_, err = db.c.Exec(
		"INSERT INTO reactions (message_id, user_id, emoji, created_at) VALUES (?, ?, ?, ?)",
		messageID, userID, emoji, time.Now(),
	)
	if err != nil {
		return ErrDatabaseError
	}

	return nil
}
