package api

import (
	"encoding/json"
	"net/http"

	"github.com/LingFengJ/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	ctx.Logger.Info("Searching for users with query: ", query)

	users, err := rt.db.SearchUsers(query)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to search users")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info("Found users: ", len(users))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode users response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
