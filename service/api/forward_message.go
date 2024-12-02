package api

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

type ForwardRequest struct {
    ConversationId string `json:"conversationId"`
}

func ForwardMessage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    messageId := vars["messageId"]

    var req ForwardRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Get original message and create new one in target conversation
    msg := forwardMessageToConversation(messageId, req.ConversationId)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(msg)
}