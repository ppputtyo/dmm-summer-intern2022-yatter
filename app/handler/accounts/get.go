package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/v1/accounts/"):]
	fmt.Println(username)
	a := h.app.Dao.Account()

	entity, err := a.FindByUsername(r.Context(), username)

	if err != nil {
		panic("not found")
	}

	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
	}
}
