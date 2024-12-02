package api

import (
    "encoding/json"
    "net/http"
)

func GetMyConversations(w http.ResponseWriter, r *http.Request) {
    // Get user ID from auth token/context
    userId := getUserFromContext(r.Context())
    
    // Get user's conversations
    conversations := getUserConversations(userId)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(conversations)
}
