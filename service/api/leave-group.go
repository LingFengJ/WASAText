package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

type LeaveGroupResponse struct {
	Success bool `json:"success"`
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Group ID required", http.StatusBadRequest)
		return
	}

	// Verify it's a group conversation
	conv, err := rt.db.GetConversation(groupID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrConversationNotFound):
			http.Error(w, "Group not found", http.StatusNotFound)
		default:
			ctx.Logger.WithError(err).Error("failed to get group")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if conv.Type != ConversationTypeGroup {
		http.Error(w, "Not a group conversation", http.StatusBadRequest)
		return
	}

	// Leave group
	err = rt.db.RemoveConversationMember(groupID, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to leave group")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := LeaveGroupResponse{Success: true}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
