package statuses

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

type PostRequest struct {
	Status    string
	media_ids []int
}

// type Post struct{
// 	ID int64
// 	Account object.Account
// 	content string
// 	create_at string

// }

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	var req PostRequest

	// json受け取る
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	// Statusの構造体
	status := new(object.Status)
	// usernameを代入
	status.Content = req.Status

	account := auth.AccountOf(r)
	status.AccountID = account.ID

	// dbを持つaccount構造体を返す
	s := h.app.Dao.Status() // domain/repository の取得
	err := s.PostStatus(r.Context(), status)
	if err != nil {
		panic("Must Implement Account Registration")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
