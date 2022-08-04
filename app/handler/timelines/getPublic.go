package timelines

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	maxID_str := r.URL.Query().Get("max_id")
	sinceID_str := r.URL.Query().Get("since_id")
	limit_str := r.URL.Query().Get("limit")

	if maxID_str == "" {
		q.MaxID = 1000000
	} else {
		res, err := strconv.Atoi(maxID_str)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}
		q.MaxID = res
	}

	if sinceID_str == "" {
		q.SinceID = -1
	} else {
		res, err := strconv.Atoi(sinceID_str)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}
		q.SinceID = res
	}

	if limit_str == "" {
		q.Limit = 40
	} else {
		res, err := strconv.Atoi(limit_str)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}
		if res > 80 {
			res = 80
		}
		q.Limit = res
	}

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
		tmp := Res{
			ID:       e.ID,
			Account:  account,
			Content:  e.Content,
			CreateAt: e.CreateAt,
		}

		res = append(res, tmp)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
