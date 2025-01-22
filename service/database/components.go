package database

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	PhotoURL   string    `json:"photoUrl,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

// Conversation represents a chat conversation (individual or group)
type Conversation struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "individual" or "group"
	Name        string    `json:"name,omitempty"`
	PhotoURL    string    `json:"photoUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	LastMessage *Message  `json:"lastMessage,omitempty"`
}

// type Conversation struct {
//     ID      string   `json:"id"`
//     Members []string `json:"members"`
//     Name    string   `json:"name"`
// }

// Message represents a chat message
type Message struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversationId"`
	SenderID       string     `json:"senderId"`
	SenderUsername string     `json:"senderUsername,omitempty"`
	Type           string     `json:"type"` // "text" or "photo"
	Content        string     `json:"content"`
	Status         string     `json:"status"` // "sent", "received", or "read"
	Timestamp      time.Time  `json:"timestamp"`
	ReplyToID      string     `json:"replyToId,omitempty"`
	Reactions      []Reaction `json:"reactions,omitempty"`
}

// Reaction represents an emoji reaction to a message
type Reaction struct {
	MessageID string    `json:"messageId"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"createdAt"`
}

// ConversationMember represents a user's membership in a conversation
type ConversationMember struct {
	ConversationID string    `json:"conversationId"`
	UserID         string    `json:"userId"`
	JoinedAt       time.Time `json:"joinedAt"`
	LastReadAt     time.Time `json:"lastReadAt"`
}

// MessageStatus tracks message delivery status per recipient
type MessageStatus struct {
	MessageID string    `json:"messageId"`
	UserID    string    `json:"userId"`
	Status    string    `json:"status"` // "received" or "read"
	UpdatedAt time.Time `json:"updatedAt"`
}
