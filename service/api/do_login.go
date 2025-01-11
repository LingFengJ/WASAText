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
		rt.baseLogger.WithError(err).Error("database error checking user existence")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var identifier string
	if exists {
		// Try to login
		identifier, err = rt.db.GetUserByCredentials(req.Name, req.Password)
		if err != nil {
			rt.baseLogger.WithError(err).Error("error getting user credentials")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	} else {
		// Create new user
		identifier, err = rt.db.CreateUser(req.Name, req.Password)
		if err != nil {
			rt.baseLogger.WithError(err).Error("error creating user")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp := LoginResponse{
		Identifier: identifier,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		rt.baseLogger.WithError(err).Error("failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
