package database

func (db *appdbimpl) UpdateLastReadTime(conversationID, userID string) error {
	query := `
        UPDATE conversation_members 
        SET last_read_at = CURRENT_TIMESTAMP 
        WHERE conversation_id = $1 AND user_id = $2`

	result, err := db.c.Exec(query, conversationID, userID)
	if err != nil {
		return ErrDatabaseError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrDatabaseError
	}

	if rowsAffected == 0 {
		return ErrConversationNotFound
	}

	return nil
}
