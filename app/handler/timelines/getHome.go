package timelines

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) GetHome(w http.ResponseWriter, r *http.Request) {
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

	myID := auth.AccountOf(r).ID

	a := h.app.Dao.Account()

	following, err := a.GetFollowing(r.Context(), myID, INF)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	userID_following := make([]int64, 0)
	for _, v := range following {
		tmp, _ := a.FindByUsername(r.Context(), v.Username)
		userID_following = append(userID_following, tmp.ID)
	}

	s := h.app.Dao.Status()
	entity, err := s.GetHomeTimelines(r.Context(), myID, query, userID_following)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	timeline := make([]Res, 0)
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
