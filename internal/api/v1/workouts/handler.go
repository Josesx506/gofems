package workouts

import (
	"database/sql"
	"encoding/json"
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

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout Workout

	// Decode the JSON body into the workout struct
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("Error decodingCreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"}) // 400
		return
	}

	createdWorkout, err := wh.store.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("Error createWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"}) // 500
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
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
		wh.logger.Printf("Error getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "invalid workout sent"}) // 404
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
		wh.logger.Printf("Error decodingUpdateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"}) // 400
		return
	}

	workout.ID = int(workoutID)

	err = wh.store.UpdateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("Error updateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update workout"}) // 500
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
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
		wh.logger.Printf("Error deleteWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to delete workout"}) // 500
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
