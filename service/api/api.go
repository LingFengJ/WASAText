package api

import (
    "github.com/gorilla/mux"
    "net/http"
)

type Handler struct {
    router *mux.Router
    // Add any dependencies like database connection here
}

func NewHandler(router *mux.Router) *Handler {
    return &Handler{
        router: router,
    }
}

func (h *Handler) RegisterRoutes() {
    // Login routes
    h.router.HandleFunc("/session", h.DoLogin).Methods(http.MethodPost)

    // User routes
    h.router.HandleFunc("/users/me/name", h.SetMyUserName).Methods(http.MethodPut)
    h.router.HandleFunc("/users/me/photo", h.SetMyPhoto).Methods(http.MethodPut)

    // Conversation routes
    h.router.HandleFunc("/conversations", h.GetMyConversations).Methods(http.MethodGet)
    h.router.HandleFunc("/conversations/{conversationId}", h.GetConversation).Methods(http.MethodGet)
    h.router.HandleFunc("/conversations/{conversationId}/messages", h.SendMessage).Methods(http.MethodPost)

    // Message routes
    h.router.HandleFunc("/messages/{messageId}", h.DeleteMessage).Methods(http.MethodDelete)
    h.router.HandleFunc("/messages/{messageId}/forward", h.ForwardMessage).Methods(http.MethodPost)
    h.router.HandleFunc("/messages/{messageId}/reactions", h.CommentMessage).Methods(http.MethodPost)
    h.router.HandleFunc("/messages/{messageId}/reactions", h.UncommentMessage).Methods(http.MethodDelete)

    // Group routes
    h.router.HandleFunc("/groups/{groupId}/leave", h.LeaveGroup).Methods(http.MethodPost)
    h.router.HandleFunc("/groups/{groupId}/members", h.AddToGroup).Methods(http.MethodPost)
    h.router.HandleFunc("/groups/{groupId}/name", h.SetGroupName).Methods(http.MethodPut)
    h.router.HandleFunc("/groups/{groupId}/photo", h.SetGroupPhoto).Methods(http.MethodPut)
}