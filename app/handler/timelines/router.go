package timelines

import (
	"net/http"
	"yatter-backend-go/app/app"

	"github.com/go-chi/chi"
)

type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()
	// r.Use(auth.Middleware(app))

	h := &handler{app: app}
	// r.With(auth.Middleware(app)).Post("/", h.Post)
	r.Get("/public", h.Get)

	return r
}
