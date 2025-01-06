package api

import (
    "encoding/json"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/LingFengJ/WASAText/service/api/reqcontext"
    "github.com/LingFengJ/WASAText/service/database"
)

type CommentMessageRequest struct {
    Emoji string `json:"emoji"`
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    messageID := ps.ByName("messageId")
    if messageID == "" {
        http.Error(w, "Message ID required", http.StatusBadRequest)
        return
    }

    var req CommentMessageRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Emoji == "" {
        http.Error(w, "Emoji required", http.StatusBadRequest)
        return
    }

    // Verify the message exists and user can access it
    msg, err := rt.db.GetMessage(messageID)
    if err != nil {
        switch err {
        case database.ErrMessageNotFound:
            http.Error(w, "Message not found", http.StatusNotFound)
        default:
            ctx.Logger.WithError(err).Error("failed to get message")
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    // Check if user is member of the conversation
    isMember, err := rt.isConversationMember(msg.ConversationID, ctx.UserID)
    if err != nil {
        ctx.Logger.WithError(err).Error("failed to check conversation membership")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    if !isMember {
        http.Error(w, "Not authorized to comment on this message", http.StatusForbidden)
        return
    }

    // Add reaction
    err = rt.db.AddReaction(messageID, ctx.UserID, req.Emoji)
    if err != nil {
        switch err {
        case database.ErrDuplicateReaction:
            http.Error(w, "Reaction already exists", http.StatusConflict)
        default:
            ctx.Logger.WithError(err).Error("failed to add reaction")
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusCreated)
}