package routes

import (
	"github.com/Josesx506/gofems/internals/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthChecker)
	return r
}
