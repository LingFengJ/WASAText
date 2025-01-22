package database

import (
	"database/sql"
)

func (db *appdbimpl) GetMessages(conversationID string, limit, offset int) ([]Message, error) {
	query := `
        SELECT 
            m.id,
            m.conversation_id,
            m.sender_id,
			u.username as sender_username,
            m.type,
            m.content,
            m.status,
            m.timestamp,
            m.reply_to_id
        FROM messages m
		LEFT JOIN users u ON m.sender_id = u.id
        WHERE m.conversation_id = $1
        ORDER BY m.timestamp DESC
        LIMIT $2 OFFSET $3`

	rows, err := db.c.Query(query, conversationID, limit, offset)
	if err != nil {
		return nil, ErrDatabaseError
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var replyToID sql.NullString      // For handling nullable reply_to_id
		var senderUsername sql.NullString // For handling nullable username

		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&senderUsername,
			&msg.Type,
			&msg.Content,
			&msg.Status,
			&msg.Timestamp,
			&replyToID,
		)
		if err != nil {
			return nil, ErrDatabaseError
		}

		if senderUsername.Valid {
			msg.SenderUsername = senderUsername.String
		}

		if replyToID.Valid {
			msg.ReplyToID = replyToID.String
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrDatabaseError
	}

	return messages, nil
}
