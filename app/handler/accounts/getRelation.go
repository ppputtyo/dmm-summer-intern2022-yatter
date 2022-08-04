package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) GetRelation(w http.ResponseWriter, r *http.Request) {
	targetUsername := r.URL.Query().Get("username")
	a := h.app.Dao.Account()

	entity, err := a.FindByUsername(r.Context(), targetUsername)

	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	targetID := entity.ID
	myID := auth.AccountOf(r).ID

	relation, err := a.GetRelation(r.Context(), myID, targetID)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(relation); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
