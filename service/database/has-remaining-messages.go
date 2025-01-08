package database

func (db *appdbimpl) HasRemainingMessages(conversationID string) (bool, error) {
	var count int
	err := db.c.QueryRow(
		"SELECT COUNT(*) FROM messages WHERE conversation_id = ?",
		conversationID,
	).Scan(&count)

	if err != nil {
		return false, ErrDatabaseError
	}

	return count > 0, nil
}
