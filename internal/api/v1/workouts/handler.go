package workouts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Josesx506/gofems/internal/utils"
)

type WorkoutHandler struct {
	store  WorkoutStore
	logger *log.Logger
}

// Accepts a WorkoutStore interface to interact with the db layer
func NewWorkoutHandler(store WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		store:  store,
		logger: logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("Error reading ID param: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	workout, err := wh.store.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("Error getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout Workout

	// Decode the JSON body into the workout struct
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid workout payload", http.StatusBadRequest) // 400
		return
	}

	createdWorkout, err := wh.store.CreateWorkout(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create workout", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("Error reading ID param: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	// Check if workout exists
	existingWorkout, err := wh.store.GetWorkoutByID(workoutID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Workout not found", http.StatusNotFound) // 404
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	// Create a new workout struct to hold the updated data
	var workout Workout

	err = json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid workout payload", http.StatusBadRequest) // 400
		return
	}

	workout.ID = int(workoutID)

	err = wh.store.UpdateWorkout(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to update workout", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("Error reading ID param: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	err = wh.store.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to delete workout", http.StatusInternalServerError) // 500
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
