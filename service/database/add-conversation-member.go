package database

import (
	"fmt"
	"time"
)

// AddConversationMember adds a new member to a conversation
func (db *appdbimpl) AddConversationMember(conversationID, userID string) error {
	now := time.Now()

	_, err := db.c.Exec(`
        INSERT INTO conversation_members (conversation_id, user_id, joined_at, last_read_at)
        VALUES ($1, $2, $3, $3)`,
		conversationID, userID, now)

	if err != nil {
		return fmt.Errorf("error adding conversation member: %w", err)
	}

	return nil
}
