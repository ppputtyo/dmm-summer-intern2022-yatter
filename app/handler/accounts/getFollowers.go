package accounts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	const INF = 10000000
	const DEFAULT_LIMIT = 40
	const MAX_LIMIT = 80

	username := chi.URLParam(r, "username")
	max_id_str := r.URL.Query().Get("max_id")
	since_id_str := r.URL.Query().Get("since_id")
	limit_str := r.URL.Query().Get("limit")

	a := h.app.Dao.Account()

	entity, err := a.FindByUsername(r.Context(), username)

	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	query := object.FollowersQuery{
		MaxID:   INF,
		SinceID: -1,
		Limit:   DEFAULT_LIMIT,
	}

	if limit_str != "" {
		query.Limit, err = strconv.Atoi(limit_str)
		if err != nil {
			httperror.BadRequest(w, err)
		}
		if query.Limit > MAX_LIMIT {
			query.Limit = MAX_LIMIT
		}
	}

	if max_id_str != "" {
		query.MaxID, err = strconv.Atoi(max_id_str)
		if err != nil {
			httperror.BadRequest(w, err)
		}
	}

	if since_id_str != "" {
		query.SinceID, err = strconv.Atoi(since_id_str)
		if err != nil {
			httperror.BadRequest(w, err)
		}
	}

	res, err := a.GetFollowers(r.Context(), entity.ID, query)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
