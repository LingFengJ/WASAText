package database

import (
	"database/sql"
)

// GetUserConversations retrieves all conversations for a user
func (db *appdbimpl) GetUserConversations(userID string) ([]Conversation, error) {
	query := `
        WITH LastMessages AS (
            SELECT 
                id,
                conversation_id,
                sender_id,
                content,
                type,
                status,
                timestamp,
                ROW_NUMBER() OVER (PARTITION BY conversation_id ORDER BY timestamp DESC) as rn
            FROM messages
        ),
        OtherMembers AS (
            -- Get other member's info for individual chats
            SELECT 
                cm.conversation_id,
                u.username as other_username,
                u.photo_url as other_photo_url
            FROM conversation_members cm
            JOIN users u ON cm.user_id = u.id
            WHERE cm.user_id != ? -- not the current user
            AND EXISTS (
                SELECT 1 FROM conversations c 
                WHERE c.id = cm.conversation_id 
                AND c.type = 'individual'
            )
        )
        SELECT 
            c.id,
            c.type,
            CASE 
                WHEN c.type = 'individual' THEN om.other_username
                ELSE c.name
            END as display_name,
            CASE 
                WHEN c.type = 'individual' THEN om.other_photo_url
                ELSE c.photo_url
            END as display_photo,
            c.created_at,
            c.modified_at,
            lm.id,
            lm.conversation_id,
            lm.sender_id,
            lm.content,
            lm.type,
            lm.status,
            lm.timestamp
        FROM conversations c
        JOIN conversation_members cm ON c.id = cm.conversation_id
        LEFT JOIN LastMessages lm ON c.id = lm.conversation_id AND lm.rn = 1
        LEFT JOIN OtherMembers om ON c.id = om.conversation_id
        WHERE cm.user_id = ?
        ORDER BY COALESCE(lm.timestamp, c.modified_at) DESC`

	rows, err := db.c.Query(query, userID, userID) // We use userID twice in the query
	if err != nil {
		return nil, ErrDatabaseError
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		// var msgContent, msgType sql.NullString
		var msgID, msgConvID, msgSenderID, msgContent, msgType, msgStatus sql.NullString
		var msgTimestamp sql.NullTime
		var displayName, displayPhoto sql.NullString

		err := rows.Scan(
			&conv.ID,
			&conv.Type,
			&displayName,
			&displayPhoto,
			&conv.CreatedAt,
			&conv.ModifiedAt,
			&msgID,
			&msgConvID,
			&msgSenderID,
			&msgContent,
			&msgType,
			&msgStatus,
			&msgTimestamp,
		)
		if err != nil {
			return nil, ErrDatabaseError
		}

		// Set the name based on what we got
		if displayName.Valid {
			conv.Name = displayName.String
		}
		if displayPhoto.Valid {
			conv.PhotoURL = displayPhoto.String
		}

		// If there's a last message, attach only preview info
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

		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrDatabaseError
	}

	return conversations, nil
}
