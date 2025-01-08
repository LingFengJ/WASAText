package database

import "time"

func (db *appdbimpl) UpdateConversation(conv *Conversation) error {
	if conv.Type != "individual" && conv.Type != "group" {
		return ErrInvalidConversationType
	}

	if conv.Type == "group" && conv.Name == "" {
		return ErrGroupNameRequired
	}

	conv.ModifiedAt = time.Now()

	result, err := db.c.Exec(`
        UPDATE conversations 
        SET 
            type = ?,
            name = ?,
            photo_url = ?,
            modified_at = ?
        WHERE id = ?`,
		conv.Type,
		conv.Name,
		conv.PhotoURL,
		conv.ModifiedAt,
		conv.ID,
	)
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

	return nil
}
