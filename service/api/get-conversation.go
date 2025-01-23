package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

const (
	ConversationTypeGroup      = "group"
	ConversationTypeIndividual = "individual"
)

type ConversationResponse struct {
	Conversation *database.Conversation `json:"conversation"`
	Messages     []database.Message     `json:"messages"`
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Conversation ID required", http.StatusBadRequest)
		return
	}

	// First, verify user is a member of this conversation
	members, err := rt.db.GetConversationMembers(conversationID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to get conversation members")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	isMember := false
	for _, member := range members {
		if member.UserID == ctx.UserID {
			isMember = true
			break
		}
	}

	if !isMember {
		http.Error(w, "Not authorized to view this conversation", http.StatusForbidden)
		return
	}

	// Get conversation details
	conversation, err := rt.db.GetConversation(conversationID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to get conversation")
		switch {
		case errors.Is(err, database.ErrConversationNotFound):
			http.Error(w, "Conversation not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Get messages (using default limit and offset)
	messages, err := rt.db.GetMessages(conversationID, 50, 0)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to get messages")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update last_read_at for the current user
	err = rt.db.UpdateLastReadTime(conversationID, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to update last read time")
		// Continue anyway as this is not critical
	}

	// Update message statuses
	for i := range messages {
		// Skip messages sent by the current user
		if messages[i].SenderID != ctx.UserID {
			err := rt.db.UpdateMessageStatus(messages[i].ID, ctx.UserID, "read")
			if err != nil {
				ctx.Logger.WithError(err).Error("failed to update message status")
				continue
			}
		}

		statuses, err := rt.db.GetMessageStatus(messages[i].ID)
		if err != nil {
			ctx.Logger.WithError(err).Error("failed to get message statuses")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			continue
		}

		if messages[i].SenderID == ctx.UserID {
			allRead := true
			allReceived := true
			receiverCount := 0

			// Count number of receivers
			for _, member := range members {
				if member.UserID == messages[i].SenderID {
					continue
				}
				receiverCount++
			}

			// Check actual statuses
			statusCount := 0
			for _, status := range statuses {
				if status.UserID == messages[i].SenderID {
					continue
				}
				statusCount++
				if status.Status != "read" {
					allRead = false
				}
				if status.Status != "received" && status.Status != "read" {
					allReceived = false
				}
			}

			var newStatus string
			if allRead && statusCount == receiverCount {
				newStatus = "read"
			} else if allReceived && statusCount == receiverCount {
				newStatus = "received"
			}

			if newStatus != "" && newStatus != messages[i].Status && messages[i].Status != "read" {
				err = rt.db.UpdateMessageAggregateStatus(messages[i].ID, newStatus)
				if err != nil {
					ctx.Logger.WithError(err).Error("failed to update message aggregate status")
				}
				messages[i].Status = newStatus
			}
		}
	}

	// // For messages not sent by current user, mark as read
	// err := rt.db.UpdateMessageStatus(messages[i].ID, ctx.UserID, "read")
	// if err != nil {
	// 	ctx.Logger.WithError(err).Error("failed to update message status")
	// 	continue
	// }

	// // For group conversations, we need to check if all members have read
	// if conversation.Type == database.ConversationTypeGroup {
	// 	allRead := true
	// 	for _, member := range members {
	// 		if member.UserID != messages[i].SenderID &&
	// 			member.LastReadAt.Before(messages[i].Timestamp) {
	// 			allRead = false
	// 			break
	// 		}
	// 	}
	// 	if allRead {
	// 		messages[i].Status = "read"
	// 	} else {
	// 		messages[i].Status = "received"
	// 	}
	// }
	// }

	// Get reactions for each message
	for i := range messages {
		reactions, err := rt.db.GetMessageReactions(messages[i].ID)
		if err != nil {
			ctx.Logger.WithError(err).Error("failed to get message reactions")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		messages[i].Reactions = reactions
	}

	response := ConversationResponse{
		Conversation: conversation,
		Messages:     messages,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
