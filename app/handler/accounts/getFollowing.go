package accounts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	const DEFAULT_LIMIT = 40
	const MAX_LIMIT = 80

	username := chi.URLParam(r, "username")
	limit_str := r.URL.Query().Get("limit")

	a := h.app.Dao.Account()

	entity, err := a.FindByUsername(r.Context(), username)

	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	limit := DEFAULT_LIMIT
	if limit_str != "" {
		limit, err = strconv.Atoi(limit_str)
		if err != nil {
			httperror.BadRequest(w, err)
		}
		if limit > MAX_LIMIT {
			limit = MAX_LIMIT
		}
	}

	res, err := a.GetFollowing(r.Context(), entity.ID, limit)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
