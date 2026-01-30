package workouts

import (
	"github.com/go-chi/chi/v5"
)

func WorkoutRouter() chi.Router {
	r := chi.NewRouter()
	handler := NewWorkoutHandler()

	r.Get("/{id}", handler.HandleGetWorkoutByID)
	r.Post("/", handler.HandleCreateWorkout)

	return r
}
