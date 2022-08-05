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

	var limit int
	if limit_str == "" {
		limit = DEFAULT_LIMIT
	} else {
		limit, err = strconv.Atoi(limit_str)
		if err != nil {
			httperror.BadRequest(w, err)
		}
	}
	if limit > MAX_LIMIT {
		limit = MAX_LIMIT
	}

	var max_id int
	if max_id_str == "" {
		max_id = INF
	} else {
		max_id, err = strconv.Atoi(max_id_str)
		if err != nil {
			httperror.BadRequest(w, err)
		}
	}

	var since_id int
	if since_id_str == "" {
		since_id = -1
	} else {
		since_id, err = strconv.Atoi(since_id_str)
		if err != nil {
			httperror.BadRequest(w, err)
		}
	}

	q := object.FollowersQuery{
		MaxID:   max_id,
		SinceID: since_id,
		Limit:   limit,
	}

	res, err := a.GetFollowers(r.Context(), entity.ID, q)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
