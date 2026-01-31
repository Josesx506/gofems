### DB Stores
Stores are part of the handlers but primarily focused on communicating with the db. It allows db operations and is essentially a replacement for ORMs with custom functions and raw SQL.

When using raw sql without ORMs, we don't define table schemas in a models directory like flask or a prisma schema file. In Go, we can define `structs` that have the same fields as the db, which makes it easier to keep track of the valid columns. The fields of the struct will have json tags to enable encoding and decoding between json payloads and struct fields with matching name tags as we use them within handlers
```go
type Workout struct {
	ID              int            `json:"id"`
	UserID          int            `json:"user_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}
```
These structs can usually be saved in a file that shares the same name as the package e.g. `workouts.go` or models.go.

We can create stores as structs that enable connections to the db
```go
type PostgresWorkoutStore struct {
	db *sql.DB
}
```
To prevent tight coupling, we can define generic interfaces that multipe db stores can utilize, as long as we define the necessary methods. For this project, we use a `WorkoutStore` interface. By defining a WorkoutStore interface with method signatures, the application layer can work with the interface without being tied to a specific database implementation, making future migrations easier.

An example method from the interface can be 
```go
func (pgStore *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
    tx, err := pgStore.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // rollback transaction if not committed
    ...
}
```
This allows us to connect to the db, execute queries, and implement transactions with rollback/commit. When writing multiple workout entries for a new workout, a ***deferred transaction rollback*** is used, and errors are checked at each stage of the process. If an error occurs, the method returns nil and the error, and the transaction can be rolled back to its previous state.

To execute a query, 
- we write the query with raw sql and placeholder for query parameters. The money symbol ($1, $2, $3, $4) is used to represent query parameters, with the number corresponding to the order of values being inserted.
- execute the query using `tx.QueryRow()` to get a single row response. This returns only an error. Use `db.Query()` to get a response for multiple rows. This returns two values, the rows as a variable, and an error. Both method takes in the raw sql as the **first argument**, then the json decoded values / query parameter values as subsequent arguments. 
- we check if the query executed successfully or returned an error.

```go
    ... // db connection
    // Step 1
	query := `
	INSERT INTO workouts (title, description, duration_minutes, calories_burned) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id
	`
    // Step 2
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)

    // Step 3
    if err != nil {
		return nil, err
	}
```

At the end of the transaction, we commit all our changes
```go
    ... // executed queries
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

    // No errors encountered
    return workout, nil
```

If any errors were encountered, the defered rollback method will be triggered, and all previous insert statements will be rendered void.

>[!NOTE]
> Because each handler requires a db connection, and the db is initialized inside the `app.Application` struct, there's quite a bit of prop drillingto handle the dependency injection. An alternative would be to define the store at the top of the main api router, move all store logic to the store package, and initialize new stores which are passed to the handlers.