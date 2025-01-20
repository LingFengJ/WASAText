package database

func (db *appdbimpl) GetMessageReactions(messageID string) ([]Reaction, error) {
	rows, err := db.c.Query(`
        SELECT 
            r.message_id,
            r.user_id,
            u.username,
            r.emoji,
            r.created_at
        FROM reactions r
        JOIN users u ON r.user_id = u.id
        WHERE r.message_id = ?`,
		messageID)

	if err != nil {
		return nil, ErrDatabaseError
	}
	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var reaction Reaction
		err := rows.Scan(
			&reaction.MessageID,
			&reaction.UserID,
			&reaction.Username,
			&reaction.Emoji,
			&reaction.CreatedAt,
		)
		if err != nil {
			return nil, ErrDatabaseError
		}
		reactions = append(reactions, reaction)
	}

	return reactions, nil
}
