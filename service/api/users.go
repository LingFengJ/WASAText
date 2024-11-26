package api

import (
    "encoding/json"
    "net/http"
)

type UserNameUpdateRequest struct {
    Name string `json:"name"`
}

func (h *Handler) SetMyUserName(w http.ResponseWriter, r *http.Request) {
    var req UserNameUpdateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // TODO: Implement username update logic
    // 1. Validate new username
    // 2. Check if username is available
    // 3. Update username in database

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"name": req.Name})
}

func (h *Handler) SetMyPhoto(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement photo upload logic
    // 1. Parse multipart form
    // 2. Validate image
    // 3. Save image and generate URL
    // 4. Update user profile
}