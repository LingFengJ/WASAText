// In get-message.go
package database

import (
	"database/sql"
)

func (db *appdbimpl) GetMessage(messageID string) (*Message, error) {
	var msg Message
	err := db.c.QueryRow(`
        SELECT 
            id,
            conversation_id,
            sender_id,
            type,
            content,
            status,
            timestamp,
            reply_to_id
        FROM messages
        WHERE id = ?`,
		messageID,
	).Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.SenderID,
		&msg.Type,
		&msg.Content,
		&msg.Status,
		&msg.Timestamp,
		&msg.ReplyToID,
	)

	if err == sql.ErrNoRows {
		return nil, ErrMessageNotFound
	}
	if err != nil {
		return nil, ErrDatabaseError
	}

	return &msg, nil
}
