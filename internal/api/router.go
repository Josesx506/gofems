package api

import (
	apiv1 "github.com/Josesx506/gofems/internal/api/v1"
	"github.com/Josesx506/gofems/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware" // middleware logger
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/health", app.HealthChecker)

	// v1 api routes
	r.Mount("/v1", apiv1.ApiV1Router(app))

	return r
}
