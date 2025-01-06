package api

import (
    "encoding/json"
    "io"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/gofrs/uuid"
    "github.com/LingFengJ/WASAText/service/api/reqcontext"
    "github.com/LingFengJ/WASAText/service/database"
    "strings"
)

type PhotoResponse struct {
    PhotoURL string `json:"photoUrl"`
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    groupID := ps.ByName("groupId")
    if groupID == "" {
        http.Error(w, "Group ID required", http.StatusBadRequest)
        return
    }

    // Check content type
    contentType := r.Header.Get("Content-Type")
    if !strings.HasPrefix(contentType, "image/") {
        http.Error(w, "Content-Type must be an image", http.StatusBadRequest)
        return
    }

    // Get current conversation
    conv, err := rt.db.GetConversation(groupID)
    if err != nil {
        switch err {
        case database.ErrConversationNotFound:
            http.Error(w, "Group not found", http.StatusNotFound)
        default:
            ctx.Logger.WithError(err).Error("failed to get group")
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    if conv.Type != "group" {
        http.Error(w, "Not a group conversation", http.StatusBadRequest)
        return
    }

    // Verify user is member
    isMember, err := rt.isConversationMember(groupID, ctx.UserID)
    if err != nil {
        ctx.Logger.WithError(err).Error("failed to check group membership")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    if !isMember {
        http.Error(w, "Not authorized to change group photo", http.StatusForbidden)
        return
    }

    // Limit file size (e.g., 5MB)
    r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

    // Read image data
    imageData, err := io.ReadAll(r.Body)
    if err != nil {
        if err.Error() == "http: request body too large" {
            http.Error(w, "Image too large", http.StatusRequestEntityTooLarge)
            return
        }
        http.Error(w, "Failed to read image data", http.StatusBadRequest)
        return
    }

    // Generate unique filename
    filename, err := uuid.NewV4()
    if err != nil {
        ctx.Logger.WithError(err).Error("failed to generate filename")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Save image and get URL
    photoURL, err := rt.db.UpdateGroupPhoto(groupID, filename.String(), imageData)
    if err != nil {
        ctx.Logger.WithError(err).Error("failed to update group photo")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Return response
    response := PhotoResponse{PhotoURL: photoURL}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}