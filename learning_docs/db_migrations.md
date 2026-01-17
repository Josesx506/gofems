### Installation
To enable migrations, we need to install the goose package globally within the application. <br>
Goose is a database migration tool for Go that allows modifications to database schema structure, 
enabling developers to add, change, or revert database changes safely and flexibly. <br>
This can be done with `go install github.com/pressly/goose/v3/cmd/goose@latest`. This is 
because goose has to be added to Go PATH. If your bashrc or bash_profile files does not 
contain a valid path for Go, add this line `export PATH=$HOME/go/bin:$PATH` and source the 
file or restart your terminal to access goose. Once installed check the version
```bash
vscode ➜ /workspaces/gofems (main)$ goose version
# Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND

# or

# Set environment key
# GOOSE_DRIVER=DRIVER
# GOOSE_DBSTRING=DBSTRING

# Usage: goose [OPTIONS] COMMAND

# Drivers:
#     postgres
#     mysql
#     sqlite3
#     mssql
#     redshift
#     tidb
#     clickhouse
#     ydb
#     turso
#  .....
```
or `goose -version` for a shorter response. You can also verify your installation by running 
`ls -l ~/go/bin | grep goose` in terminal. Additional info can be found in the official goose 
installation [guide](https://pressly.github.io/goose/installation/).<br>

Finally, install goose as a module, so it can be imported into files `go get github.com/pressly/goose/v3/cmd/goose@latest`.

### Setup
Create a migrations module, and include a file that allows inclusion of sql files directly into 
our compiled application binary.
```go
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
```
This code uses the `embed` package (introduced in Go 1.16) to bundle your SQL migration files directly 
into your compiled application binary. In simple terms: it treats your external .sql files as if they 
were variables inside your code.
- `//go:embed *.sql`: This is a compiler directive. It tells the Go compiler: "Find every file ending 
    in .sql in this current directory and prepare them to be embedded."
- `var FS embed.FS`: This creates a read-only virtual file system (embed.FS) containing all those matched 
    SQL files. 

Go migration libraries (like `goose`, `golang-migrate`, or `sql-migrate`) can read directly from an embed.FS.

***Important***<br>
- The directive `//go:embed` must be placed **directly** above the variable it's populating (no empty lines).
- The files must be in the **same directory** (or a sub-directory) as the `.go` file containing the directive.
- The variable must be at the **package level** (not inside a function).

### Running migrations
Manually create the migration sql files in the `migrations` directory, and then create a `Migrate()` function 
inside the [database setup module](/internals/store/database.go). For each sql command, the goose down command 
must be the inverse of the up command.
- If `Up` creates a table, `Down` should drop that table.
- If `Up` adds a column, `Down` should remove that column.
- If `Up` inserts data, `Down` should delete that data.

After setting up the db store Open() method, run the migration up command as part of the application setup in 
the [App module](/internals/app/app.go). Now when you start the application
```bash
vscode ➜ /workspaces/gofems (main) $ go run main.go 
Connected to database....
2026/01/17 08:43:33 OK   00001_users.sql (2.97ms)
2026/01/17 08:43:33 OK   00002_workouts.sql (1.01ms)
2026/01/17 08:43:33 OK   00003_workout_entries.sql (409.5µs)
2026/01/17 08:43:33 goose: successfully migrated database to version: 3
2026/01/17 08:43:33 We are running our api on port 8080
```
you should see the migrated version.

### COnnecting the DB to the store structs
Post migration, we can define our DB types as structs, similar to how an ORM will have a variable for accessing 
the db. In Go, we can include struct tags to improve json conversion. <br>
A struct tag helps with encoding and decoding JSON into a struct, allowing easy conversion between JSON payloads 
and struct fields with matching tag names.
```go
type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}
```
Fileds with pointer values explicitly allow checking for nil values, which provides more flexibility in handling 
optional or potentially unset numeric fields.