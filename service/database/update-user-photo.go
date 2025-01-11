package database

import (
	"os"
	"path"
	"path/filepath"
	"time"
)

func (db *appdbimpl) UpdateUserPhoto(userID string, filename string, imageData []byte) (string, error) {
	// Create photos directory if it doesn't exist
	photosDir := "data/photos"
	if err := os.MkdirAll(photosDir, 0755); err != nil {
		return "", ErrDatabaseError
	}

	// Create file with extension
	photoPath := filepath.Join(photosDir, filename+".jpg")
	if err := os.WriteFile(photoPath, imageData, 0644); err != nil {
		return "", ErrDatabaseError
	}

	// Create URL for the photo
	// photoURL := fmt.Sprintf("/photos/%s.jpg", filename)
	photoURL := path.Join("/photos", filename+".jpg")

	// Update database with photo URL
	result, err := db.c.Exec(`
        UPDATE users 
        SET photo_url = ?, modified_at = ?
        WHERE id = ?`,
		photoURL, time.Now(), userID,
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
		return "", ErrUserNotFound
	}

	return photoURL, nil
}
