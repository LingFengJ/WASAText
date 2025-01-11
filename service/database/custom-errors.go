package database

import "errors"

var (
	// User-related errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameTaken      = errors.New("username is already taken")
	ErrInvalidUsername    = errors.New("invalid username format")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Conversation-related errors
	ErrConversationNotFound    = errors.New("conversation not found")
	ErrNotConversationMember   = errors.New("user is not a member of this conversation")
	ErrInvalidConversationType = errors.New("invalid conversation type")
	ErrGroupNameRequired       = errors.New("group name is required")

	// Message-related errors
	ErrMessageNotFound       = errors.New("message not found")
	ErrInvalidMessageType    = errors.New("invalid message type")
	ErrMessageNotOwner       = errors.New("user is not the message owner")
	ErrInvalidMessageContent = errors.New("invalid message content")
	ErrInvalidMessageStatus  = errors.New("invalid message status")

	// Reaction-related errors
	ErrInvalidEmoji      = errors.New("invalid emoji")
	ErrReactionNotFound  = errors.New("reaction not found")
	ErrDuplicateReaction = errors.New("user has already reacted with this emoji")

	// Permission-related errors
	ErrNotGroupAdmin    = errors.New("user is not a group admin")
	ErrCannotLeaveGroup = errors.New("cannot leave group: user is the last member")

	// General errors
	ErrInvalidID      = errors.New("invalid identifier format")
	ErrDatabaseError  = errors.New("database error")
	ErrTableCreation  = errors.New("error creating database table")
	ErrTableCheck     = errors.New("error checking table existence")
	ErrInvalidRequest = errors.New("invalid request")
)

// ValidationError represents an error that occurs during data validation
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// DatabaseError wraps database-specific errors with additional context
type DatabaseError struct {
	Op  string // Operation being performed
	Err error  // Original error
}

func (e *DatabaseError) Error() string {
	return e.Op + ": " + e.Err.Error()
}

// NewDatabaseError creates a new DatabaseError
func NewDatabaseError(op string, err error) *DatabaseError {
	return &DatabaseError{
		Op:  op,
		Err: err,
	}
}
