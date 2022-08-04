package accounts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	limit := r.URL.Query().Get("limit")
	a := h.app.Dao.Account()

	entity, err := a.FindByUsername(r.Context(), username)

	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	var l int
	if limit == "" {
		l = 40
	} else {
		l, err = strconv.Atoi(limit)
		if err != nil {
			httperror.BadRequest(w, err)
		}
	}

	if l > 80 {
		l = 80
	}

	res, err := a.GetFollowing(r.Context(), entity.ID, l)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
