package api

import (
	"encoding/json"
	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type SetGroupNameRequest struct {
	Name string `json:"name"`
}

// type SetGroupNameResponse struct {
//     Name string `json:"name"`
// }

type SetGroupNameResponse SetGroupNameRequest // Define response type as alias of request type

// Convert request to response
func (req SetGroupNameRequest) ToResponse() SetGroupNameResponse {
	return SetGroupNameResponse(req)
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Group ID required", http.StatusBadRequest)
		return
	}

	var req SetGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Group name required", http.StatusBadRequest)
		return
	}

	// Get current conversation
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

	// Verify user is member
	isMember, err := rt.isConversationMember(groupID, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to check group membership")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Not authorized to rename this group", http.StatusForbidden)
		return
	}

	// Update name
	conv.Name = req.Name
	err = rt.db.UpdateConversation(conv)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to update group name")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	// response := SetGroupNameResponse{Name: req.Name}
	response := req.ToResponse()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
