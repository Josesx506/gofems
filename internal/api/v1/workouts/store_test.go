package workouts

import (
	"testing"

	"github.com/Josesx506/gofems/internal/store"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWorkout(t *testing.T) {
	db := store.SetupTestDB(t, "../../../../migrations/")
	defer db.Close()

	pgStore := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "Valid Workout",
			workout: &Workout{
				Title:           "Morning Routine",
				Description:     "A quick morning workout",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench Press",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(135.5),
						Notes:        "warm up properly",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Workout with invalid entries",
			workout: &Workout{
				Title:           "full body workout",
				Description:     "A full body workout session",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         4,
						Reps:         IntPtr(60),
						Notes:        "keep form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "squats",
						Sets:            3,
						Reps:            IntPtr(12),
						DurationSeconds: IntPtr(60), // Invalid: both Reps and DurationSeconds set
						Weight:          FloatPtr(185.0),
						Notes:           "full depth",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := pgStore.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			// Retrieve the workout from the database to verify entries
			retrievedWorkout, err := pgStore.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retrievedWorkout.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrievedWorkout.Entries))

			for i, entry := range retrievedWorkout.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, entry.ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, entry.Sets)
				assert.Equal(t, tt.workout.Entries[i].Reps, entry.Reps)
				assert.Equal(t, tt.workout.Entries[i].Weight, entry.Weight)
				assert.Equal(t, tt.workout.Entries[i].Notes, entry.Notes)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, entry.OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
