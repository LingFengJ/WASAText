package api

import (
    "encoding/json"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/LingFengJ/WASAText/service/api/reqcontext"
    "github.com/LingFengJ/WASAText/service/database"
)

type DeleteMessageResponse struct {
    Success bool `json:"success"`
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    messageID := ps.ByName("messageId")
    if messageID == "" {
        http.Error(w, "Message ID required", http.StatusBadRequest)
        return
    }

    // Verify user owns the message
    message, err := rt.db.GetMessage(messageID)
    if err != nil {
        ctx.Logger.WithError(err).Error("failed to get message")
        switch err {
        case database.ErrMessageNotFound:
            http.Error(w, "Message not found", http.StatusNotFound)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    if message.SenderID != ctx.UserID {
        http.Error(w, "Not authorized to delete this message", http.StatusForbidden)
        return
    }

    // Delete message
    err = rt.db.DeleteMessage(messageID)
    if err != nil {
        ctx.Logger.WithError(err).Error("failed to delete message")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Return success response
    response := DeleteMessageResponse{Success: true}
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        ctx.Logger.WithError(err).Error("failed to encode response")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}