package timelines

import (
	"encoding/json"
	"fmt"
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

	// onry_media, _ := strconv.Atoi(r.URL.Query().Get("only_media"))
	// max_id, err := strconv.Atoi(r.URL.Query().Get("max_id"))
	// if err != nil {
	// 	max_id = -1
	// }
	// since_id, _ := strconv.Atoi(r.URL.Query().Get("since_id"))
	// limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	// if err != nil {
	// 	limit = 40
	// }
	// if limit > 80 {
	// 	limit = 80
	// }

	fmt.Println(q)

	s := h.app.Dao.Status()

	entity, err := s.GetPublicTimelines(r.Context(), q)

	if err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
	}
}
