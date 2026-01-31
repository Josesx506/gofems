package workouts

import (
	"github.com/Josesx506/gofems/internal/app"
	"github.com/go-chi/chi/v5"
)

func WorkoutRouter(app *app.Application) chi.Router {

	// Initialize router, handler, and store
	r := chi.NewRouter()
	// Store requires global db connection
	store := &PostgresWorkoutStore{db: app.DB}
	handler := NewWorkoutHandler(store, app.Logger)

	// Define subroutes
	r.Get("/{id}", handler.HandleGetWorkoutByID)
	r.Put("/{id}", handler.HandleUpdateWorkoutByID)
	r.Post("/", handler.HandleCreateWorkout)
	r.Delete("/{id}", handler.HandleDeleteWorkoutByID)

	return r
}
