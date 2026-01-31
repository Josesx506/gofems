package workouts

import (
	"database/sql"
)

// DB connector struct
type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

// Use interface to decouple DB (postgres) from application layer
// Any other DB e.g. mysql,mongo etc. can implement this interface
type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
	// GetWorkoutOwner(id int64) (int, error)
}

// Define methods for PostgresWorkoutStore to implement WorkoutStore interface
func (pgStore *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pgStore.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // rollback transaction if not committed

	// Insert workout
	query := `
	INSERT INTO workouts (title, description, duration_minutes, calories_burned) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id
	`
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes,
		workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}

	// Insert each workout entry as a row
	for _, entry := range workout.Entries {
		entryQuery := `
		INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
		`
		// Uses the returned workout id from the previous insert and scans returned entry id
		err = tx.QueryRow(entryQuery, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps,
			entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pgStore *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	workout := &Workout{} // Initialize empty workout

	query := `
	SELECT id, title, description, duration_minutes, calories_burned
	FROM workouts
	WHERE id = $1
	`
	err := pgStore.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description,
		&workout.DurationMinutes, &workout.CaloriesBurned)

	if err == sql.ErrNoRows {
		return nil, nil // No workout found
	}

	if err != nil {
		return nil, err
	}

	entriesQuery := `
	SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
	FROM workout_entries
	WHERE workout_id = $1
	ORDER BY order_index ASC
	`
	rows, err := pgStore.db.Query(entriesQuery, workout.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after processing

	for rows.Next() {
		entry := WorkoutEntry{}
		err := rows.Scan(&entry.ID, &entry.ExerciseName, &entry.Sets, &entry.Reps,
			&entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (pgStore *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pgStore.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateQuery := `
	UPDATE workouts
	SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
	WHERE id = $5
	`
	result, err := tx.Exec(updateQuery, workout.Title, workout.Description,
		workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Delete existing workout entries
	_, err = tx.Exec(`DELETE FROM workout_entries WHERE workout_id = $1`, workout.ID)
	if err != nil {
		return err
	}

	// Insert updated workout entries
	for _, entry := range workout.Entries {
		entryQuery := `
		INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		// Uses the returned workout id from the previous insert and scans returned entry id
		_, err := tx.Exec(entryQuery, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps,
			entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pgStore *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	result, err := pgStore.db.Exec(`DELETE FROM workouts WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
