package database

func (db *appdbimpl) DeleteMessage(messageID string) error {
    // Begin transaction
    tx, err := db.c.Begin()
    if err != nil {
        return ErrDatabaseError
    }
    defer tx.Rollback()

    // Delete message status entries first due to foreign key constraint
    _, err = tx.Exec("DELETE FROM message_status WHERE message_id = ?", messageID)
    if err != nil {
        return ErrDatabaseError
    }

    // Delete reactions
    _, err = tx.Exec("DELETE FROM reactions WHERE message_id = ?", messageID)
    if err != nil {
        return ErrDatabaseError
    }

    // Delete the message
    result, err := tx.Exec("DELETE FROM messages WHERE id = ?", messageID)
    if err != nil {
        return ErrDatabaseError
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return ErrDatabaseError
    }
    if rows == 0 {
        return ErrMessageNotFound
    }

    return tx.Commit()
}