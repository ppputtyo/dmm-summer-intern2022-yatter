package timelines

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

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

	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
