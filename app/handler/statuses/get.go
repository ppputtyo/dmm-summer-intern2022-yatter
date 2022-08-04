package statuses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	statusID := chi.URLParam(r, "statusID")

	fmt.Println(statusID)
	s := h.app.Dao.Status()

	statusID_int64, _ := strconv.ParseInt(statusID, 10, 64)

	entity, err := s.FindByPostID(r.Context(), statusID_int64)

	if err != nil {
		panic(err)
	}

	a := h.app.Dao.Account() // domain/repository の取得
	account, err := a.FindByID(r.Context(), entity.AccountID)

	if err != nil {
		panic(err)
	}

	res := Res{}
	res.ID = entity.ID
	res.Account = account
	res.Content = entity.Content
	res.CreateAt = entity.CreateAt

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
	}
}
