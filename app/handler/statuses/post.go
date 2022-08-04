package statuses

import (
	"encoding/json"
	"fmt"
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

type Res struct {
	ID       int64           `json:"id"`
	Account  *object.Account `json:"account"`
	Content  string          `json:"content"`
	CreateAt object.DateTime `json:"create_at"`
}

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
		httperror.InternalServerError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")

	status, err = s.FindByID(r.Context(), status.ID)
	if err != nil {
		panic(err)
	}

	res := Res{}
	res.ID = status.ID
	res.Account = account
	res.Content = status.Content
	res.CreateAt = status.CreateAt

	fmt.Println(res)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
