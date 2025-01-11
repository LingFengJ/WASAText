package database

func (db *appdbimpl) DeleteConversation(conversationID string) error {
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

	// Delete all reactions to messages in this conversation
	_, err = tx.Exec(`
        DELETE FROM reactions 
        WHERE message_id IN (
            SELECT id FROM messages 
            WHERE conversation_id = ?
        )`, conversationID)
	if err != nil {
		return ErrDatabaseError
	}

	// Delete all message status entries
	_, err = tx.Exec(`
        DELETE FROM message_status 
        WHERE message_id IN (
            SELECT id FROM messages 
            WHERE conversation_id = ?
        )`, conversationID)
	if err != nil {
		return ErrDatabaseError
	}

	// Delete all messages
	_, err = tx.Exec(
		"DELETE FROM messages WHERE conversation_id = ?",
		conversationID)
	if err != nil {
		return ErrDatabaseError
	}

	// Delete all conversation members
	_, err = tx.Exec(
		"DELETE FROM conversation_members WHERE conversation_id = ?",
		conversationID)
	if err != nil {
		return ErrDatabaseError
	}

	// Delete the conversation itself
	result, err := tx.Exec(
		"DELETE FROM conversations WHERE id = ?",
		conversationID)
	if err != nil {
		return ErrDatabaseError
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return ErrDatabaseError
	}
	if rows == 0 {
		return ErrConversationNotFound
	}

	return tx.Commit()
}
