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

func (h *handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	const INF = 10000000
	const DEFAULT_LIMIT = 40
	const MAX_LIMIT = 80

	query := object.GetTimelineQuery{
		OnlyMedia: r.URL.Query().Get("only_media"),
		MaxID:     INF,
		SinceID:   -1,
		Limit:     DEFAULT_LIMIT,
	}

	maxID_str := r.URL.Query().Get("max_id")
	sinceID_str := r.URL.Query().Get("since_id")
	limit_str := r.URL.Query().Get("limit")

	if maxID_str != "" {
		res, err := strconv.Atoi(maxID_str)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}
		query.MaxID = res
	}

	if sinceID_str != "" {
		res, err := strconv.Atoi(sinceID_str)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}
		query.SinceID = res
	}

	if limit_str != "" {
		res, err := strconv.Atoi(limit_str)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}
		if res > MAX_LIMIT {
			res = MAX_LIMIT
		}
		query.Limit = res
	}

	s := h.app.Dao.Status()

	entity, err := s.GetPublicTimelines(r.Context(), query)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	timeline := make([]Res, 0)

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

		timeline = append(timeline, tmp)
	}

	if err := json.NewEncoder(w).Encode(timeline); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
