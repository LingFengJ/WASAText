package api

import (
	"encoding/json"
	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type UpdateUsernameRequest struct {
	Name string `json:"name"`
}

//	type UpdateUsernameResponse struct {
//	    Name string `json:"name"`
//	}
type UpdateUsernameResponse UpdateUsernameRequest // Define response type as alias of request type

// Convert request to response
func (req UpdateUsernameRequest) ToResponse() UpdateUsernameResponse {
	return UpdateUsernameResponse(req)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req UpdateUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate username
	if len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, "Username must be between 3 and 16 characters", http.StatusBadRequest)
		return
	}

	// Update username in database
	err := rt.db.UpdateUsername(ctx.UserID, req.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to update username")

		switch err {
		case database.ErrUsernameTaken:
			http.Error(w, "Username already taken", http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// // Send response
	// resp := UpdateUsernameResponse{
	//     Name: req.Name,
	// }
	resp := req.ToResponse()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
