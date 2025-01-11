package database

import (
	"errors"
)

func (db *appdbimpl) RemoveConversationMember(conversationID, userID string) error {
	result, err := db.c.Exec(`
        DELETE FROM conversation_members 
        WHERE conversation_id = $1 AND user_id = $2`,
		conversationID, userID)

	if err != nil {
		var (
			ErrRemovingConversationMember = errors.New("failed to remove conversation member")
		)
		return ErrRemovingConversationMember
	}

	rows, err := result.RowsAffected()
	if err != nil {
		var (
			ErrCheckingrows = errors.New("error checking affected rows")
		)
		return ErrCheckingrows
	}

	if rows == 0 {
		return ErrNotConversationMember
	}

	return nil
}
