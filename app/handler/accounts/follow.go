package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	targetUsername := chi.URLParam(r, "username")
	fmt.Println(targetUsername)
	a := h.app.Dao.Account()

	entity, err := a.FindByUsername(r.Context(), targetUsername)

	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	targetID := entity.ID
	myID := auth.AccountOf(r).ID

	err = a.Follow(r.Context(), myID, targetID)

	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	relation, err := a.GetRelation(r.Context(), myID, targetID)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(relation); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
