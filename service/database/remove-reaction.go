package database

func (db *appdbimpl) RemoveReaction(messageID string, userID string) error {
	result, err := db.c.Exec(
		"DELETE FROM reactions WHERE message_id = ? AND user_id = ?",
		messageID, userID,
	)
	if err != nil {
		return ErrDatabaseError
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return ErrDatabaseError
	}
	if rows == 0 {
		return ErrReactionNotFound
	}

	return nil
}
