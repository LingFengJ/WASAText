// In send-message.go
package api

import (
    "encoding/json"
    "net/http"
    "time"
    "github.com/julienschmidt/httprouter"
    "github.com/gofrs/uuid"
    "github.com/LingFengJ/WASAText/service/api/reqcontext"
    "github.com/LingFengJ/WASAText/service/database"
)

// In send-message.go
type SendMessageRequest struct {
    Content       string   `json:"content"`                  // Required for all messages
    Type         string   `json:"type"`                     // "text" or "photo"
    RecipientName string   `json:"recipientName,omitempty"` // Required only for first individual message
    GroupName    string   `json:"groupName,omitempty"`      // Optional, only for group creation
    Members      []string `json:"members,omitempty"`        // Optional, only for group creation
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    conversationID := ps.ByName("conversationId")

    // Parse request body
    var req SendMessageRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate message
    if req.Content == "" {
        http.Error(w, "Message content required", http.StatusBadRequest)
        return
    }
    if req.Type != "text" && req.Type != "photo" {
        http.Error(w, "Invalid message type", http.StatusBadRequest)
        return
    }

    // Check if this is a new conversation
    if conversationID == "" {
        var err error
        // Check if it's a group creation request
        if req.GroupName != "" {
            if len(req.Members) == 0 {
                http.Error(w, "Members required for group creation", http.StatusBadRequest)
                return
            }
            conversationID, err = rt.createNewGroup(ctx.UserID, req.GroupName, req.Members)
        } else {
            // Original individual chat logic
            if req.RecipientName == "" {
                http.Error(w, "Recipient name required for new conversation", http.StatusBadRequest)
                return
            }
            conversationID, err = rt.createNewIndividualChat(ctx.UserID, req.RecipientName)
        }
        if err != nil {
            switch err {
            case database.ErrUserNotFound:
                http.Error(w, "Recipient not found", http.StatusNotFound)
            default:
                ctx.Logger.WithError(err).Error("failed to create conversation")
                http.Error(w, "Internal server error", http.StatusInternalServerError)
            }
            return
        }
    }else {
        // For existing conversations, verify user is still a member
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

    // Create message
    message := &database.Message{
        ConversationID: conversationID,
        SenderID:      ctx.UserID,
        Type:          req.Type,
        Content:       req.Content,
        Status:        "sent",
        Timestamp:     time.Now(),
    }

    if err := rt.db.CreateMessage(message); err != nil {
        ctx.Logger.WithError(err).Error("failed to create message")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(message)
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
        Type:       "individual",
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
        Type:       "group",
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