package timelines

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

type Res struct {
	ID       int64           `json:"id"`
	Account  *object.Account `json:"account"`
	Content  string          `json:"content"`
	CreateAt object.DateTime `json:"create_at"`
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	q := object.Query{}

	q.OnlyMedia = r.URL.Query().Get("only_media")
	q.MaxID = r.URL.Query().Get("max_id")
	q.SinceID = r.URL.Query().Get("since_id")
	q.Limit = r.URL.Query().Get("limit")

	s := h.app.Dao.Status()

	entity, err := s.GetPublicTimelines(r.Context(), q)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res := make([]Res, 0)

	a := h.app.Dao.Account() // domain/repository の取得

	for _, e := range entity {
		account, err := a.FindByID(r.Context(), e.AccountID)

		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
		tmp := Res{}
		tmp.ID = e.ID
		tmp.Account = account
		tmp.Content = e.Content
		tmp.CreateAt = e.CreateAt

		res = append(res, tmp)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
