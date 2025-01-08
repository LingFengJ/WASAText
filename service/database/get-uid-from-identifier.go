package database

func (db *appdbimpl) GetUserIDFromIdentifier(identifier string) (string, error) {
	var id string
	err := db.c.QueryRow(
		"SELECT id FROM users WHERE identifier = ?",
		identifier,
	).Scan(&id)

	if err != nil {
		return "", ErrUserNotFound
	}

	return id, nil
}
