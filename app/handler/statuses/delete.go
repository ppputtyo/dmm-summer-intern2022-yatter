package statuses

import (
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	statusID := chi.URLParam(r, "statusID")
	s := h.app.Dao.Status()

	statusID_int64, err := strconv.ParseInt(statusID, 10, 64)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	account := auth.AccountOf(r)

	entity, err := s.FindByID(r.Context(), statusID_int64)

	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if entity.AccountID != account.ID {
		httperror.BadRequest(w, err)
		return
	}

	if err := s.DeleteStatus(r.Context(), statusID_int64); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
