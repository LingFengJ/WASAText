package api

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

func LeaveGroup(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    groupId := vars["groupId"]
    
    // Get user ID from auth token/context
    userId := getUserFromContext(r.Context())
    
    // Remove user from group
    success := removeUserFromGroup(userId, groupId)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"success": success})
}