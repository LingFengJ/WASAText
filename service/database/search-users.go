package database

import (
	"database/sql"
)

func (db *appdbimpl) SearchUsers(query string) ([]User, error) {
	var rows *sql.Rows
	var err error

	if query == "" {
		// If no query provided, return all users
		rows, err = db.c.Query(`
            SELECT 
                id,
                username,
                COALESCE(photo_url, '') as photo_url,
                created_at,
                modified_at
            FROM users 
            ORDER BY username ASC
            LIMIT 50`) // Reasonable limit for all users
	} else {
		// Search with query
		rows, err = db.c.Query(`
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
			query+"%")
	}

	if err != nil {
		return nil, ErrDatabaseError
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var photoUrl sql.NullString

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
