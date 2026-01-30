package api

import (
	apiv1 "github.com/Josesx506/gofems/internals/api/v1"
	"github.com/Josesx506/gofems/internals/app"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5/middleware" // middleware logger
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// r.Use(middleware.Logger)

	r.Get("/health", app.HealthChecker)

	// v1 api routes
	r.Mount("/api/v1", apiv1.ApiV1Router())

	return r
}
