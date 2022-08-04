package accounts

import (
	"net/http"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	// v1/accounts/
	r.Post("/", h.Create)
	r.Get("/{username}", h.Get)
	r.Get("/{username}/following", h.GetFollowing)
	r.Get("/{username}/followers", h.GetFollowers)
	r.With(auth.Middleware(app)).Get("/relationships", h.GetRelation)
	r.With(auth.Middleware(app)).Post("/{username}/follow", h.Follow)

	return r
}
