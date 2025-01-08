package database

func (db *appdbimpl) GetMessageStatus(messageID string) ([]MessageStatus, error) {
	// Check if message exists
	var exists bool
	err := db.c.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM messages WHERE id = ?)",
		messageID,
	).Scan(&exists)
	if err != nil {
		return nil, ErrDatabaseError
	}
	if !exists {
		return nil, ErrMessageNotFound
	}

	// Get all status entries for the message
	rows, err := db.c.Query(`
        SELECT 
            message_id,
            user_id,
            status,
            updated_at
        FROM message_status
        WHERE message_id = ?
        ORDER BY updated_at DESC`,
		messageID,
	)
	if err != nil {
		return nil, ErrDatabaseError
	}
	defer rows.Close()

	var statuses []MessageStatus
	for rows.Next() {
		var status MessageStatus
		err := rows.Scan(
			&status.MessageID,
			&status.UserID,
			&status.Status,
			&status.UpdatedAt,
		)
		if err != nil {
			return nil, ErrDatabaseError
		}
		statuses = append(statuses, status)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrDatabaseError
	}

	return statuses, nil
}
