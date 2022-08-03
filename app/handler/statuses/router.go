package statuses

import (
	"net/http"
	"yatter-backend-go/app/app"

	"github.com/go-chi/chi"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	// h := &handler{app: app}
	// // v1/accounts/
	// r.Post("/", h.Create)
	// r.Get("/{usename}", h.Get)

	return r
}