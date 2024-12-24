package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handles APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login endpoints (no auth required)
	h := rt
	h.router.POST("/session", h.doLogin)

	// User endpoints
	h.router.PUT("/users/me/name", h.wrap(h.setMyUserName))
	h.router.PUT("/users/me/photo", h.wrap(h.setMyPhoto))

	// Conversation endpoints
	h.router.GET("/conversations", h.wrap(h.getMyConversations))
	h.router.GET("/conversations/:conversationId", h.wrap(h.getConversation))
	h.router.POST("/conversations/:conversationId/messages", h.wrap(h.sendMessage))

	// Message endpoints
	h.router.DELETE("/messages/:messageId", h.wrap(h.deleteMessage))
	h.router.POST("/messages/:messageId/forward", h.wrap(h.forwardMessage))
	h.router.POST("/messages/:messageId/reactions", h.wrap(h.commentMessage))
	h.router.DELETE("/messages/:messageId/reactions", h.wrap(h.uncommentMessage))

	// Group endpoints
	h.router.POST("/groups/:groupId/leave", h.wrap(h.leaveGroup))
	h.router.POST("/groups/:groupId/members", h.wrap(h.addToGroup))
	h.router.PUT("/groups/:groupId/name", h.wrap(h.setGroupName))
	h.router.PUT("/groups/:groupId/photo", h.wrap(h.setGroupPhoto))

	// Special endpoints
	h.router.GET("/liveness", h.liveness)

	return h.router
}

// // Custom type for context keys to avoid collisions
// type contextKey string

// const UserIDKey = contextKey("userId")

// // wrap handles authentication and user identification before passing to the actual handler
// func (h *Handler) wrap(handler httprouter.Handle) httprouter.Handle {
//     return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//         // Get the Authorization header
//         authHeader := r.Header.Get("Authorization")
//         if authHeader == "" {
//             http.Error(w, "Authorization header required", http.StatusUnauthorized)
//             return
//         }

//         // Check for Bearer prefix
//         parts := strings.Split(authHeader, " ")
//         if len(parts) != 2 || parts[0] != "Bearer" {
//             http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
//             return
//         }

//         // Get the token
//         token := parts[1]
//         if token == "" {
//             http.Error(w, "Token required", http.StatusUnauthorized)
//             return
//         }

//         // Add user ID to context
//         ctx := r.Context()
//         ctx = context.WithValue(ctx, UserIDKey, token)
//         r = r.WithContext(ctx)

//         // Call the actual handler
//         handler(w, r, ps)  //// maybe add ctx as a parameter following fantastic coffee
//     }
// }


// // liveness is a basic health check endpoint
// func (h *_router) liveness(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	w.WriteHeader(http.StatusOK)
// }


