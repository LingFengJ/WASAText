package database

import (
	"database/sql"
)

func (db *appdbimpl) GetConversation(conversationID string) (*Conversation, error) {
	query := `
        SELECT 
            c.id,
            c.type,
            c.name,
            c.photo_url,
            c.created_at,
            c.modified_at,
            m.id as last_message_id,
            m.conversation_id as last_message_conv_id,
            m.sender_id as last_message_sender_id,
            m.type as last_message_type,
            m.content as last_message_content,
            m.status as last_message_status,
            m.timestamp as last_message_time
        FROM conversations c
        LEFT JOIN (
            SELECT *
            FROM messages
            WHERE (conversation_id, timestamp) IN (
                SELECT conversation_id, MAX(timestamp)
                FROM messages
                GROUP BY conversation_id
            )
        ) m ON c.id = m.conversation_id
        WHERE c.id = ?`

	var conv Conversation
	var msgID, msgConvID, msgSenderID, msgType, msgContent, msgStatus sql.NullString
	var msgTimestamp sql.NullTime

	err := db.c.QueryRow(query, conversationID).Scan(
		&conv.ID,
		&conv.Type,
		&conv.Name,
		&conv.PhotoURL,
		&conv.CreatedAt,
		&conv.ModifiedAt,
		&msgID,
		&msgConvID,
		&msgSenderID,
		&msgType,
		&msgContent,
		&msgStatus,
		&msgTimestamp,
	)

	if err == sql.ErrNoRows {
		return nil, ErrConversationNotFound
	}
	if err != nil {
		return nil, ErrDatabaseError
	}

	// If there's a last message, attach it with all fields
	if msgContent.Valid && msgType.Valid && msgTimestamp.Valid {
		conv.LastMessage = &Message{
			ID:             msgID.String,
			ConversationID: msgConvID.String,
			SenderID:       msgSenderID.String,
			Type:           msgType.String,
			Content:        msgContent.String,
			Status:         msgStatus.String,
			Timestamp:      msgTimestamp.Time,
		}
	}

	return &conv, nil
}
