/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		// return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
)

type AppDatabase interface {
	Ping() error

	// User management
	CheckUserExists(username string) (bool, error)
	CreateUser(username string, password string) (string, error)           // Returns identifier
	GetUserByCredentials(username string, password string) (string, error) // Returns identifier
	GetUserIDFromIdentifier(identifier string) (string, error)
	GetUserIDByUsername(username string) (string, error)
	SearchUsers(query string) ([]User, error) // Returns a list of users

	// User profile management
	UpdateUsername(userID string, newUsername string) error
	UpdateUserPhoto(userID string, filename string, imageData []byte) (string, error)

	// // Conversation management
	GetUserConversations(userID string) ([]Conversation, error)
	GetConversation(conversationID string) (*Conversation, error)
	CreateConversation(conv *Conversation, members []string) error
	UpdateConversation(conv *Conversation) error
	DeleteConversation(conversationID string) error
	HasRemainingMessages(conversationID string) (bool, error)
	UpdateGroupPhoto(conversationID, filename string, imageData []byte) (string, error)

	// // Conversation members
	AddConversationMember(conversationID, userID string) error
	RemoveConversationMember(conversationID, userID string) error
	GetConversationMembers(conversationID string) ([]ConversationMember, error)

	// // Message management
	GetMessages(conversationID string, limit, offset int) ([]Message, error)
	GetMessage(messageID string) (*Message, error)
	CreateMessage(msg *Message) error
	UpdateMessage(msg *Message) error
	DeleteMessage(messageID string) error
	AddReaction(messageID, userID, emoji string) error
	RemoveReaction(messageID, userID string) error

	// // Message status
	UpdateMessageStatus(messageID, userID, status string) error
	GetMessageStatus(messageID string) ([]MessageStatus, error)
}

type appdbimpl struct {
	c *sql.DB
}

// Table creation statements
const (
	usersTableCreationStatement = `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        identifier TEXT UNIQUE NOT NULL,
        photo_url TEXT,
        created_at TIMESTAMP NOT NULL,
        modified_at TIMESTAMP NOT NULL
    );`

	conversationsTableCreationStatement = `
    CREATE TABLE IF NOT EXISTS conversations (
        id TEXT PRIMARY KEY,
        type TEXT NOT NULL CHECK (type IN ('individual', 'group')),
        name TEXT,
        photo_url TEXT,
        created_at TIMESTAMP NOT NULL,
        modified_at TIMESTAMP NOT NULL
    );`

	conversationMembersTableCreationStatement = `
    CREATE TABLE IF NOT EXISTS conversation_members (
        conversation_id TEXT,
        user_id TEXT,
        joined_at TIMESTAMP NOT NULL,
        last_read_at TIMESTAMP NOT NULL,
        PRIMARY KEY (conversation_id, user_id),
        FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

	messagesTableCreationStatement = `
    CREATE TABLE IF NOT EXISTS messages (
        id TEXT PRIMARY KEY,
        conversation_id TEXT NOT NULL,
        sender_id TEXT NOT NULL,
        type TEXT NOT NULL CHECK (type IN ('text', 'photo')),
        content TEXT NOT NULL,
        status TEXT NOT NULL CHECK (status IN ('sent', 'received', 'read')),
        timestamp TIMESTAMP NOT NULL,
        reply_to_id TEXT,
        FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
        FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (reply_to_id) REFERENCES messages(id) ON DELETE SET NULL
    );`

	messageStatusTableCreationStatement = `
    CREATE TABLE IF NOT EXISTS message_status (
        message_id TEXT,
        user_id TEXT,
        status TEXT NOT NULL CHECK (status IN ('received', 'read')),
        updated_at TIMESTAMP NOT NULL,
        PRIMARY KEY (message_id, user_id),
        FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

	reactionsTableCreationStatement = `
    CREATE TABLE IF NOT EXISTS reactions (
        message_id TEXT,
        user_id TEXT,
        emoji TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        PRIMARY KEY (message_id, user_id),
        FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`
)

func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Map of table names to their creation statements
	tableMapping := map[string]string{
		"users":                usersTableCreationStatement,
		"conversations":        conversationsTableCreationStatement,
		"conversation_members": conversationMembersTableCreationStatement,
		"messages":             messagesTableCreationStatement,
		"message_status":       messageStatusTableCreationStatement,
		"reactions":            reactionsTableCreationStatement,
	}

	// Check if each table exists. If not, create it
	for tableName, sqlStmt := range tableMapping {
		err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, tableName).Scan(&tableName)

		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.Exec(sqlStmt)
			if err != nil {
				// return nil, fmt.Errorf("error creating database structure for table %s: %w", tableName, err)
				return nil, ErrTableCreation
			}
		} else if err != nil {
			// return nil, fmt.Errorf("error checking table %s existence: %w", tableName, err)
			return nil, ErrTableCheck
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
