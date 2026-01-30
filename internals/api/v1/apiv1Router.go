package apiv1

import (
	"github.com/Josesx506/gofems/internals/api/v1/workouts"
	"github.com/go-chi/chi/v5"
)

func ApiV1Router() chi.Router {
	r := chi.NewRouter()
	// Handler doesn't need to be imported since they're in the same package
	v1Handler := NewApiV1Handler()

	r.Get("/health", v1Handler.Health)

	r.Mount("/workouts", workouts.WorkoutRouter())

	return r
}
