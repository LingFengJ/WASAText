package api

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

type MessageRequest struct {
    Content string `json:"content"`
    Type    string `json:"type"`
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    conversationId := vars["conversationId"]

    var req MessageRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Create and save the message
    msg := createMessage(conversationId, req.Content, req.Type)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(msg)
}