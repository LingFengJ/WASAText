package database

import (
	"fmt"
)

func (db *appdbimpl) GetConversationMembers(conversationID string) ([]ConversationMember, error) {
    rows, err := db.c.Query(`
        SELECT conversation_id, user_id, joined_at, last_read_at
        FROM conversation_members
        WHERE conversation_id = $1`,
        conversationID)
    
    if err != nil {
        return nil, fmt.Errorf("error getting conversation members: %w", err)
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
            return nil, fmt.Errorf("error scanning conversation member: %w", err)
        }
        members = append(members, member)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating conversation members: %w", err)
    }
    
    return members, nil
}