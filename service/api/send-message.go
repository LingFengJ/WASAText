// In send-message.go
package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// In send-message.go
type SendMessageRequest struct {
	Content       string   `json:"content"`                 // Required for all messages
	Type          string   `json:"type"`                    // "text" or "photo"
	RecipientName string   `json:"recipientName,omitempty"` // Required only for first individual message
	GroupName     string   `json:"groupName,omitempty"`     // Optional, only for group creation
	Members       []string `json:"members,omitempty"`       // Optional, only for group creation
	ReplyToID     string   `json:"replyToId,omitempty"`     // Optional, only for replies
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID := ps.ByName("conversationId")

	var req SendMessageRequest

	// Handle different content types
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse form data
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("content")
		if err != nil {
			http.Error(w, "Content file required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			ctx.Logger.WithError(err).Error("failed to read file")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Convert to base64
		base64Content := base64.StdEncoding.EncodeToString(fileBytes)

		req = SendMessageRequest{
			Content:       base64Content,
			Type:          "photo",
			ReplyToID:     r.FormValue("replyToId"),
			RecipientName: r.FormValue("recipientName"),
			GroupName:     r.FormValue("groupName"),
		}

		if members := r.Form["members[]"]; len(members) > 0 {
			req.Members = members
		}
	} else {
		// Handle JSON request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// If it's a photo type message, validate base64 content
		if req.Type == "photo" {
			// Verify that the content is valid base64
			if _, err := base64.StdEncoding.DecodeString(req.Content); err != nil {
				http.Error(w, "Invalid base64 image content", http.StatusBadRequest)
				return
			}
		}
	}

	// Validate message content
	if req.Content == "" {
		http.Error(w, "Message content required", http.StatusBadRequest)
		return
	}
	if req.Type != "text" && req.Type != "photo" {
		http.Error(w, "Invalid message type", http.StatusBadRequest)
		return
	}

	// Handle new conversation creation
	if conversationID == "" {
		var err error
		if req.GroupName != "" {
			// Group chat creation
			if len(req.Members) == 0 {
				http.Error(w, "Members required for group creation", http.StatusBadRequest)
				return
			}
			conversationID, err = rt.createNewGroup(ctx.UserID, req.GroupName, req.Members)
		} else {
			// Individual chat creation
			if req.RecipientName == "" {
				http.Error(w, "Recipient name required for new conversation", http.StatusBadRequest)
				return
			}
			conversationID, err = rt.createNewIndividualChat(ctx.UserID, req.RecipientName)
		}
		if err != nil {
			switch {
			case errors.Is(err, database.ErrUserNotFound):
				http.Error(w, "Recipient not found", http.StatusNotFound)
			default:
				ctx.Logger.WithError(err).Error("failed to create conversation")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
	} else {
		// Verify membership for existing conversation
		isMember, err := rt.isConversationMember(conversationID, ctx.UserID)
		if err != nil {
			ctx.Logger.WithError(err).Error("failed to check conversation membership")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !isMember {
			http.Error(w, "Not authorized to send messages to this conversation", http.StatusForbidden)
			return
		}
	}

	username, err := rt.db.GetUsernameByIdentifier(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to get username")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create message
	message := &database.Message{
		ConversationID: conversationID,
		SenderID:       ctx.UserID,
		SenderUsername: username,
		Type:           req.Type,
		Content:        req.Content,
		Status:         "sent",
		Timestamp:      time.Now(),
		ReplyToID:      req.ReplyToID,
	}

	if err := rt.db.CreateMessage(message); err != nil {
		ctx.Logger.WithError(err).Error("failed to create message")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Helper for creating individual chats
func (rt *_router) createNewIndividualChat(senderID, recipientName string) (string, error) {
	// Get recipient's ID
	recipientID, err := rt.db.GetUserIDByUsername(recipientName)
	if err != nil {
		return "", err
	}

	convID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	conv := &database.Conversation{
		ID:         convID.String(),
		Type:       ConversationTypeIndividual,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	err = rt.db.CreateConversation(conv, []string{senderID, recipientID})
	if err != nil {
		return "", err
	}

	return conv.ID, nil
}

// Helper for creating group chats
func (rt *_router) createNewGroup(creatorID, groupName string, memberUsernames []string) (string, error) {
	// Convert usernames to IDs
	memberIDs := []string{creatorID} // Creator is automatically a member
	for _, username := range memberUsernames {
		memberID, err := rt.db.GetUserIDByUsername(username)
		if err != nil {
			return "", err
		}
		memberIDs = append(memberIDs, memberID)
	}

	convID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	conv := &database.Conversation{
		ID:         convID.String(),
		Type:       ConversationTypeGroup,
		Name:       groupName,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	err = rt.db.CreateConversation(conv, memberIDs)
	if err != nil {
		return "", err
	}

	return conv.ID, nil
}

// Helper function to check conversation membership
func (rt *_router) isConversationMember(conversationID, userID string) (bool, error) {
	members, err := rt.db.GetConversationMembers(conversationID)
	if err != nil {
		return false, err
	}

	for _, member := range members {
		if member.UserID == userID {
			return true, nil
		}
	}

	return false, nil
}
