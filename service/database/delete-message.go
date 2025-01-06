package database

import (
    "database/sql"
)

func (db *appdbimpl) DeleteMessage(messageID string) error {
    // Begin transaction
    tx, err := db.c.Begin()
    if err != nil {
        return ErrDatabaseError
    }
    defer tx.Rollback()

    // Get conversation ID before deleting the message
    var conversationID string
    err = tx.QueryRow("SELECT conversation_id FROM messages WHERE id = ?", messageID).Scan(&conversationID)
    if err == sql.ErrNoRows {
        return ErrMessageNotFound
    }
    if err != nil {
        return ErrDatabaseError
    }

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

    // Check if this was the last message in the conversation
    var messageCount int
    err = tx.QueryRow(
        "SELECT COUNT(*) FROM messages WHERE conversation_id = ?",
        conversationID,
    ).Scan(&messageCount)
    if err != nil {
        return ErrDatabaseError
    }

    // If no messages left, delete the conversation
    if messageCount == 0 {
        // Delete conversation members
        _, err = tx.Exec(
            "DELETE FROM conversation_members WHERE conversation_id = ?",
            conversationID,
        )
        if err != nil {
            return ErrDatabaseError
        }

        // Delete the conversation
        _, err = tx.Exec(
            "DELETE FROM conversations WHERE id = ?",
            conversationID,
        )
        if err != nil {
            return ErrDatabaseError
        }
    }

    return tx.Commit()
}