package database

// import "time"

func (db *appdbimpl) UpdateMessage(msg *Message) error {
    // First check if message exists and user is the sender
    var senderID string
    err := db.c.QueryRow(
        "SELECT sender_id FROM messages WHERE id = ?",
        msg.ID,
    ).Scan(&senderID)
    
    if err != nil {
        return ErrMessageNotFound
    }
    
    if senderID != msg.SenderID {
        return ErrMessageNotOwner
    }

    // Update message
    res, err := db.c.Exec(`
        UPDATE messages 
        SET content = ?,
            type = ?,
            status = ?,
            reply_to_id = ?
        WHERE id = ?`,
        msg.Content,
        msg.Type,
        msg.Status,
        msg.ReplyToID,
        msg.ID,
    )
    if err != nil {
        return ErrDatabaseError
    }

    rows, err := res.RowsAffected()
    if err != nil {
        return ErrDatabaseError
    }
    if rows == 0 {
        return ErrMessageNotFound
    }

    return nil
}