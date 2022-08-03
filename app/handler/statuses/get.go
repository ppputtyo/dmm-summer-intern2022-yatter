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

	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
	}
}
