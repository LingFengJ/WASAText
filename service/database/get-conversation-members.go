package database

import (
	"errors"
)

func (db *appdbimpl) GetConversationMembers(conversationID string) ([]ConversationMember, error) {
	rows, err := db.c.Query(`
        SELECT conversation_id, user_id, joined_at, last_read_at
        FROM conversation_members
        WHERE conversation_id = $1`,
		conversationID)

	if err != nil {
		var (
			ErrGettingConversationMember = errors.New("error in db for getting conversation member")
		)
		return nil, ErrGettingConversationMember
	}
	defer rows.Close()

	var members []ConversationMember
	for rows.Next() {
		var member ConversationMember
		err := rows.Scan(
			&member.ConversationID,
			&member.UserID,
			&member.JoinedAt,
			&member.LastReadAt,
		)
		if err != nil {
			var (
				ErrScanningConversationMember = errors.New("DB Error in scanning conversation member")
			)
			return nil, ErrScanningConversationMember
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		var (
			ErrIteratingConversationMember = errors.New("error iterating conversation members")
		)
		return nil, ErrIteratingConversationMember
	}

	return members, nil
}
