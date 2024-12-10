package api

import (
    "encoding/json"
    "net/http"
    "github.com/julienschmidt/httprouter"
)

type LoginRequest struct {
    Name     string `json:"name"`
    Password string `json:"password"` // Added password field
}

type LoginResponse struct {
    Identifier string `json:"identifier"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Check if user exists
    exists, err := rt.db.CheckUserExists(req.Name)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    var identifier string
    if exists {
        // Try to login
        identifier, err = rt.db.GetUserByCredentials(req.Name, req.Password)
        if err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }
    } else {
        // Create new user
        identifier, err = rt.db.CreateUser(req.Name, req.Password)
        if err != nil {
            http.Error(w, "Could not create user", http.StatusInternalServerError)
            return
        }
    }

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