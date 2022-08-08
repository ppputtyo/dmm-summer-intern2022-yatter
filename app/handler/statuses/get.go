package statuses

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	statusID := chi.URLParam(r, "statusID")

	s := h.app.Dao.Status()

	statusID_int64, err := strconv.ParseInt(statusID, 10, 64)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	entity, err := s.FindByID(r.Context(), statusID_int64)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	a := h.app.Dao.Account() // domain/repository の取得

	account, err := a.FindByID(r.Context(), entity.AccountID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res := Res{
		ID:       entity.ID,
		Account:  account,
		Content:  entity.Content,
		CreateAt: entity.CreateAt,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
