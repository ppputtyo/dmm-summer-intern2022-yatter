package timelines

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) GetHome(w http.ResponseWriter, r *http.Request) {
	q := object.GetTimelineQuery{}

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
	myID := auth.AccountOf(r).ID

	a := h.app.Dao.Account()
	following, _ := a.GetFollowing(r.Context(), myID, 1000000)

	userID_following := make([]int64, 0)
	for _, v := range following {
		tmp, _ := a.FindByUsername(r.Context(), v.Username)
		userID_following = append(userID_following, tmp.ID)
	}
	fmt.Println(userID_following)

	entity, err := s.GetHomeTimelines(r.Context(), myID, q, userID_following)

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res := make([]Res, 0)

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
