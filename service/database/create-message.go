package database

import (
	"time"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) CreateMessage(msg *Message) error {
	// Generate message ID if not provided
	if msg.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return ErrDatabaseError
		}
		msg.ID = id.String()
	}

	// Set timestamp if not provided
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	// Validate message type
	if msg.Type != "text" && msg.Type != "photo" {
		return ErrInvalidMessageType
	}

	// Begin transaction
	tx, err := db.c.Begin()
	if err != nil {
		return ErrDatabaseError
	}

	defer func() {
		if rerr := tx.Rollback(); rerr != nil {
			err = ErrDatabaseError
		}
	}()

	// Insert message
	_, err = tx.Exec(`
        INSERT INTO messages (
            id, conversation_id, sender_id, type, content, 
            status, timestamp, reply_to_id
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.ConversationID, msg.SenderID, msg.Type,
		msg.Content, msg.Status, msg.Timestamp, msg.ReplyToID,
	)
	if err != nil {
		return ErrDatabaseError
	}

	// Get all conversation members except sender for message status
	rows, err := tx.Query(`
        SELECT user_id 
        FROM conversation_members 
        WHERE conversation_id = ? AND user_id != ?`,
		msg.ConversationID, msg.SenderID,
	)
	if err != nil {
		return ErrDatabaseError
	}
	defer rows.Close()

	// Insert initial message status for all recipients
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return ErrDatabaseError
		}

		_, err = tx.Exec(`
            INSERT INTO message_status (
                message_id, user_id, status, updated_at
            ) VALUES (?, ?, 'sent', ?)`,
			msg.ID, userID, time.Now(),
		)
		if err != nil {
			return ErrDatabaseError
		}
	}

	if err = rows.Err(); err != nil {
		return ErrDatabaseError
	}

	return tx.Commit()
}
