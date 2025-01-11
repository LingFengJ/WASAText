// In start-conversation.go
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

type StartConversationRequest struct {
	RecipientName string `json:"recipientName"`       // Username of the person to chat with
	Type          string `json:"type"`                // "individual" or "group"
	GroupName     string `json:"groupName,omitempty"` // Required for group chats
}

func (rt *_router) startConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req StartConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate new conversation ID
	convID, err := uuid.NewV4()
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to generate conversation ID")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get recipient's ID from username
	recipientID, err := rt.db.GetUserIDByUsername(req.RecipientName)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrUserNotFound):
			http.Error(w, "Recipient not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Create conversation
	conv := &database.Conversation{
		ID:         convID.String(),
		Type:       req.Type,
		Name:       req.GroupName,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	members := []string{ctx.UserID, recipientID}
	err = rt.db.CreateConversation(conv, members)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to create conversation")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
