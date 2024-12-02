package api

import (
    "encoding/json"
    "net/http"
)

type LoginRequest struct {
    Name string `json:"name"`
}

type LoginResponse struct {
    Identifier string `json:"identifier"`
}

func DoLogin(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Simple user check/creation logic
    // In a real app, you'd want to store this in a database
    existingUser := checkUserExists(req.Name)
    if existingUser {
        // User exists, generate identifier
        identifier := generateIdentifier(req.Name)
        sendLoginResponse(w, identifier)
        return
    }

    // New user, create and generate identifier
    identifier := createNewUser(req.Name)
    sendLoginResponse(w, identifier)
}

func sendLoginResponse(w http.ResponseWriter, identifier string) {
    resp := LoginResponse{
        Identifier: identifier,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resp)
}