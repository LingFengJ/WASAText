package api

import (
	"encoding/json"
	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strings"
)

type UpdatePhotoResponse struct {
	PhotoURL string `json:"photoUrl"`
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check content type
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "Content-Type must be an image", http.StatusBadRequest)
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
	photoURL, err := rt.db.UpdateUserPhoto(ctx.UserID, filename.String(), imageData)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to update user photo")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send response
	resp := UpdatePhotoResponse{
		PhotoURL: photoURL,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
