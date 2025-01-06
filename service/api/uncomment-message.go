package api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    messageID := ps.ByName("messageId")
    if messageID == "" {
        http.Error(w, "Message ID required", http.StatusBadRequest)
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
        http.Error(w, "Not authorized to remove reactions from this message", http.StatusForbidden)
        return
    }

    // Remove reaction
    err = rt.db.RemoveReaction(messageID, ctx.UserID)
    if err != nil {
        switch err {
        case database.ErrReactionNotFound:
            http.Error(w, "Reaction not found", http.StatusNotFound)
        default:
            ctx.Logger.WithError(err).Error("failed to remove reaction")
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
}