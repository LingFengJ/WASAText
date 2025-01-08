package api

import (
	"encoding/json"
	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type AddToGroupRequest struct {
	Username string `json:"username"` // Changed from userId to username
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Group ID required", http.StatusBadRequest)
		return
	}

	var req AddToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from username
	userToAddID, err := rt.db.GetUserIDByUsername(req.Username)
	if err != nil {
		switch err {
		case database.ErrUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
		default:
			ctx.Logger.WithError(err).Error("failed to get user ID")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Verify it's a group conversation
	conv, err := rt.db.GetConversation(groupID)
	if err != nil {
		switch err {
		case database.ErrConversationNotFound:
			http.Error(w, "Group not found", http.StatusNotFound)
		default:
			ctx.Logger.WithError(err).Error("failed to get group")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if conv.Type != "group" {
		http.Error(w, "Not a group conversation", http.StatusBadRequest)
		return
	}

	// Verify current user is a member
	isMember, err := rt.isConversationMember(groupID, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to check group membership")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Not authorized to add members to this group", http.StatusForbidden)
		return
	}

	// Add new member
	err = rt.db.AddConversationMember(groupID, userToAddID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to add member to group")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
