package api

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

type GroupNameRequest struct {
    Name string `json:"name"`
}

func SetGroupName(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    groupId := vars["groupId"]

    var req GroupNameRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Update group name
    success := updateGroupName(groupId, req.Name)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"name": req.Name})
}