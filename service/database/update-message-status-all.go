package database

func (db *appdbimpl) UpdateMessageAggregateStatus(messageID string, status string) error {
	// No need to check sender - status can be updated by the system
	query := `
        UPDATE messages 
        SET status = ?
        WHERE id = ?`

	result, err := db.c.Exec(query, status, messageID)
	if err != nil {
		return ErrDatabaseError
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return ErrDatabaseError
	}
	if rows == 0 {
		return ErrMessageNotFound
	}

	return nil
}
