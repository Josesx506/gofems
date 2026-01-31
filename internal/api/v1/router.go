package apiv1

import (
	"github.com/Josesx506/gofems/internal/api/v1/workouts"
	"github.com/Josesx506/gofems/internal/app"
	"github.com/go-chi/chi/v5"
)

func ApiV1Router(app *app.Application) chi.Router {
	r := chi.NewRouter()
	// Handler doesn't need to be imported since they're in the same package
	v1Handler := NewApiV1Handler(app)

	r.Get("/health", v1Handler.Health)

	r.Mount("/workouts", workouts.WorkoutRouter(app))

	return r
}
