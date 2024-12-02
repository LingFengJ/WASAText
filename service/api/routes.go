package api

import (
    "github.com/gorilla/mux"
    "net/http"
)

func RegisterRoutes(router *mux.Router) {
    // Login routes
    router.HandleFunc("/session", DoLogin).Methods(http.MethodPost)

    // Message routes
    router.HandleFunc("/conversations/{conversationId}/messages", SendMessage).Methods(http.MethodPost)
    router.HandleFunc("/messages/{messageId}/forward", ForwardMessage).Methods(http.MethodPost)
    
    // Group routes
    router.HandleFunc("/groups/{groupId}/leave", LeaveGroup).Methods(http.MethodPost)
    
    // Add other routes similarly...
}