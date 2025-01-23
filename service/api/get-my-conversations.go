package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from context
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get conversations from database
	conversations, err := rt.db.GetUserConversations(userID)
	if err != nil {
		// Log the error
		ctx.Logger.WithError(err).Error("failed to get user conversations")

		// Check for specific errors
		switch {
		case errors.Is(err, database.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound)
		case errors.Is(err, database.ErrDatabaseError):
			http.Error(w, "Database error", http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// If no conversations, return empty array instead of null
	if conversations == nil {
		conversations = []database.Conversation{}
	}

	// Update status to "received" for all unread messages
	for _, conv := range conversations {
		messages, err := rt.db.GetMessages(conv.ID, 50, 0) // Get all recent messages
		if err != nil {
			ctx.Logger.WithError(err).Error("failed to get messages")
			continue
		}

		// Get conversation members for group status check
		members, err := rt.db.GetConversationMembers(conv.ID)
		if err != nil {
			ctx.Logger.WithError(err).Error("failed to get conversation members")
			continue
		}

		for _, msg := range messages {
			if msg.SenderID == ctx.UserID {
				continue
			}

			// Update status to "received" for any unseen message
			if msg.Status == "sent" {
				err = rt.db.UpdateMessageStatus(msg.ID, ctx.UserID, "received")
				if err != nil {
					ctx.Logger.WithError(err).Error("failed to update message status")
					continue
				}
			}

			// Get all statuses for this message
			statuses, err := rt.db.GetMessageStatus(msg.ID)
			if err != nil {
				ctx.Logger.WithError(err).Error("failed to get message statuses")
				continue
			}

			// Calculate aggregate status
			receiverCount := 0
			for _, member := range members {
				if member.UserID == msg.SenderID {
					continue
				}
				receiverCount++
			}

			allReceived := true
			for _, status := range statuses {
				if status.UserID == msg.SenderID {
					continue
				}
				if status.Status != "received" && status.Status != "read" {
					allReceived = false
					break
				}
			}
			// Update the message's aggregate status if needed
			if allReceived && msg.Status == "sent" && msg.Status != "read" {
				err = rt.db.UpdateMessageAggregateStatus(msg.ID, "received")
				if err != nil {
					ctx.Logger.WithError(err).Error("failed to update message aggregate status")
				}
			}
		}
	}

	// for _, msg := range messages {
	// 	if msg.SenderID != ctx.UserID && msg.Status == "sent" {
	// 		if conv.Type == database.ConversationTypeGroup {
	// 			// check if all members have received for group
	// 			allReceived := true
	// 			for _, member := range members {
	// 				if member.UserID != ctx.UserID &&
	// 				   member.LastReadAt.Before(msg.Timestamp){
	// 					allReceived = false
	// 					break
	// 				   }
	// 			}

	// 			if allReceived {
	// 				err := rt.db.UpdateMessageStatus(msg.ID, ctx.UserID, "received")
	// 				if err != nil {
	// 					ctx.Logger.WithError(err).Error("failed to update message status")
	// 				}
	// 			}
	// 		} else{
	// 			// For individual conversations, update status to "received"
	// 			err := rt.db.UpdateMessageStatus(msg.ID, ctx.UserID, "received")
	// 			if err != nil {
	// 				ctx.Logger.WithError(err).Error("failed to update message status")
	// 			}
	// 		}
	// 	}
	// }
	// }

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode conversations response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
