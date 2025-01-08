package database

import (
	"github.com/gofrs/uuid"
	"time"
)

func (db *appdbimpl) CreateConversation(conv *Conversation, members []string) error {
	if conv.Type != "individual" && conv.Type != "group" {
		return ErrInvalidConversationType
	}

	if conv.Type == "group" && conv.Name == "" {
		return ErrGroupNameRequired
	}

	// Begin transaction
	tx, err := db.c.Begin()
	if err != nil {
		return ErrDatabaseError
	}
	defer tx.Rollback()

	// Generate ID if not provided
	if conv.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return ErrDatabaseError
		}
		conv.ID = id.String()
	}

	now := time.Now()
	conv.CreatedAt = now
	conv.ModifiedAt = now

	// Create conversation
	_, err = tx.Exec(`
        INSERT INTO conversations (
            id, type, name, photo_url, created_at, modified_at
        ) VALUES (?, ?, ?, ?, ?, ?)`,
		conv.ID, conv.Type, conv.Name, conv.PhotoURL, conv.CreatedAt, conv.ModifiedAt,
	)
	if err != nil {
		return ErrDatabaseError
	}

	// Add members
	for _, memberID := range members {
		_, err = tx.Exec(`
            INSERT INTO conversation_members (
                conversation_id, user_id, joined_at, last_read_at
            ) VALUES (?, ?, ?, ?)`,
			conv.ID, memberID, now, now,
		)
		if err != nil {
			return ErrDatabaseError
		}
	}

	return tx.Commit()
}
