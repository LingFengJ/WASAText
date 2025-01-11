package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from context
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get conversations from database
	conversations, err := rt.db.GetUserConversations(userID)
	if err != nil {
		// Log the error
		ctx.Logger.WithError(err).Error("failed to get user conversations")

		// Check for specific errors
		switch {
		case errors.Is(err, database.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound)
		case errors.Is(err, database.ErrDatabaseError):
			http.Error(w, "Database error", http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// If no conversations, return empty array instead of null
	if conversations == nil {
		conversations = []database.Conversation{}
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode conversations response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
