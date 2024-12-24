package database

import (
	"fmt"
)

func (db *appdbimpl) RemoveConversationMember(conversationID, userID string) error {
    result, err := db.c.Exec(`
        DELETE FROM conversation_members 
        WHERE conversation_id = $1 AND user_id = $2`,
        conversationID, userID)
    
    if err != nil {
        return fmt.Errorf("error removing conversation member: %w", err)
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error checking affected rows: %w", err)
    }
    
    if rows == 0 {
        return ErrNotConversationMember
    }
    
    return nil
}