package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"
)

func (db *appdbimpl) UpdateGroupPhoto(groupID string, filename string, imageData []byte) (string, error) {
	// First verify it's a group
	var convType string
	err := db.c.QueryRow("SELECT type FROM conversations WHERE id = ?", groupID).Scan(&convType)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrConversationNotFound
	}
	if err != nil {
		return "", ErrDatabaseError
	}
	if convType != ConversationTypeGroup {
		return "", ErrInvalidConversationType
	}

	// Create photos directory if it doesn't exist
	photosDir := "data/photos/groups"
	if err := os.MkdirAll(photosDir, 0755); err != nil {
		return "", ErrDatabaseError
	}

	// Create file with extension
	photoPath := filepath.Join(photosDir, filename+".jpg")
	if err := os.WriteFile(photoPath, imageData, 0644); err != nil {
		return "", ErrDatabaseError
	}

	// Create URL for the photo
	photoURL := filepath.Join("/photos/groups", filename+".jpg")

	// Update database with photo URL and modified time
	result, err := db.c.Exec(`
        UPDATE conversations 
        SET photo_url = ?, modified_at = ?
        WHERE id = ?`,
		photoURL, time.Now(), groupID,
	)
	if err != nil {
		// Clean up file if database update fails
		os.Remove(photoPath)
		return "", ErrDatabaseError
	}

	rows, err := result.RowsAffected()
	if err != nil {
		os.Remove(photoPath)
		return "", ErrDatabaseError
	}
	if rows == 0 {
		os.Remove(photoPath)
		return "", ErrConversationNotFound
	}

	return photoURL, nil
}
