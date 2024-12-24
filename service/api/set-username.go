package api

import (
    "encoding/json"
    "net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/LingFengJ/WASAText/service/api/reqcontext"
    // "github.com/LingFengJ/WASAText/service/database"
)

type UserNameUpdateRequest struct {
    Name string `json:"name"`
}


func (h *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
    var req UserNameUpdateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // TODO: Implement username update logic
    // 1. Validate new username
    // 2. Check if username is available
    // 3. Update username in database

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"name": req.Name})
}