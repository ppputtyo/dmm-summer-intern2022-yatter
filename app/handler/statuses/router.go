package statuses

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
	// r.Use(auth.Middleware(app))

	h := &handler{app: app}
	r.With(auth.Middleware(app)).Post("/", h.Post)

	return r
}
