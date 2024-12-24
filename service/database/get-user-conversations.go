package database 

import (
	"database/sql"
	"fmt"
)
// GetUserConversations retrieves all conversations for a user
func (db *appdbimpl) GetUserConversations(userID string) ([]Conversation, error) {
    query := `
        WITH LastMessages AS (
            SELECT 
                m.conversation_id,
                m.content,
                m.type,
                m.timestamp,
                ROW_NUMBER() OVER (PARTITION BY m.conversation_id ORDER BY m.timestamp DESC) as rn
            FROM messages m
        )
        SELECT 
            c.id,
            c.type,
            c.name,
            c.photo_url,
            c.created_at,
            c.modified_at,
            lm.content,
            lm.type,
            lm.timestamp
        FROM conversations c
        JOIN conversation_members cm ON c.id = cm.conversation_id
        LEFT JOIN LastMessages lm ON c.id = lm.conversation_id AND lm.rn = 1
        WHERE cm.user_id = $1
        ORDER BY COALESCE(lm.timestamp, c.modified_at) DESC`

    rows, err := db.c.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("error querying conversations: %w", err)
    }
    defer rows.Close()

    var conversations []Conversation
    for rows.Next() {
        var conv Conversation
        var msgContent, msgType sql.NullString
        var msgTimestamp sql.NullTime

        err := rows.Scan(
            &conv.ID,
            &conv.Type,
            &conv.Name,
            &conv.PhotoURL,
            &conv.CreatedAt,
            &conv.ModifiedAt,
            &msgContent,
            &msgType,
            &msgTimestamp,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning conversation: %w", err)
        }

        // If there's a last message, attach it to the conversation
        if msgContent.Valid && msgType.Valid && msgTimestamp.Valid {
            conv.LastMessage = &Message{
                ConversationID: conv.ID,
                Content:       msgContent.String,
                Type:         msgType.String,
                Timestamp:    msgTimestamp.Time,
            }
        }

        conversations = append(conversations, conv)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating conversations: %w", err)
    }

    return conversations, nil
}