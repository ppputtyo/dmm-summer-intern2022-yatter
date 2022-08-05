package accounts

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/accounts`
type AddRequest struct {
	Username string
	Password string
}

// Handle request for `POST /v1/accounts`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req AddRequest

	// json受け取る
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	// Accountの構造体
	account := new(object.Account)
	// usernameを代入
	account.Username = req.Username

	// passwordを代入
	if err := account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// dbを持つaccount構造体を返す
	a := h.app.Dao.Account() // domain/repository の取得
	err := a.CreateNewAccount(r.Context(), *account)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	entity, err := a.FindByUsername(r.Context(), account.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
