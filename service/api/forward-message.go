package api

import (
    "encoding/json"
    "net/http"
    "time"
    "github.com/julienschmidt/httprouter"
    "github.com/LingFengJ/WASAText/service/api/reqcontext"
    "github.com/LingFengJ/WASAText/service/database"
)

type ForwardMessageRequest struct {
    ConversationID string `json:"conversationId"`
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    messageID := ps.ByName("messageId")
    if messageID == "" {
        http.Error(w, "Message ID required", http.StatusBadRequest)
        return
    }

    // Parse request
    var req ForwardMessageRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Get original message
    originalMessage, err := rt.db.GetMessage(messageID)
    if err != nil {
        ctx.Logger.WithError(err).Error("failed to get original message")
        switch err {
        case database.ErrMessageNotFound:
            http.Error(w, "Message not found", http.StatusNotFound)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    // Verify user can access target conversation
    members, err := rt.db.GetConversationMembers(req.ConversationID)
    if err != nil {
        ctx.Logger.WithError(err).Error("failed to get conversation members")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    isMember := false
    for _, member := range members {
        if member.UserID == ctx.UserID {
            isMember = true
            break
        }
    }

    if !isMember {
        http.Error(w, "Not authorized to forward message to this conversation", http.StatusForbidden)
        return
    }

    // Create forwarded message
    newMessage := &database.Message{
        ConversationID: req.ConversationID,
        SenderID:      ctx.UserID,
        Type:          originalMessage.Type,
        Content:       originalMessage.Content,
        Status:        "sent",
        Timestamp:     time.Now(),
    }

    err = rt.db.CreateMessage(newMessage)
    if err != nil {
        ctx.Logger.WithError(err).Error("failed to create forwarded message")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Return new message
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(newMessage); err != nil {
        ctx.Logger.WithError(err).Error("failed to encode response")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}