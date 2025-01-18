package database

import (
	"database/sql"
)

func (db *appdbimpl) SearchUsers(query string) ([]User, error) {
	// Query with nullable fields
	rows, err := db.c.Query(`
        SELECT 
            id,
            username,
            COALESCE(photo_url, '') as photo_url,
            created_at,
            modified_at
        FROM users 
        WHERE username LIKE ?
        ORDER BY username ASC
        LIMIT 10`,
		"%"+query+"%")

	if err != nil {
		return nil, ErrDatabaseError
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var photoUrl sql.NullString // Handle nullable photo_url

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&photoUrl,
			&user.CreatedAt,
			&user.ModifiedAt,
		)
		if err != nil {
			return nil, ErrDatabaseError
		}

		// Handle the nullable photo_url
		if photoUrl.Valid {
			user.PhotoURL = photoUrl.String
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrDatabaseError
	}

	if users == nil {
		users = []User{} // Return empty slice instead of nil
	}

	return users, nil
}
