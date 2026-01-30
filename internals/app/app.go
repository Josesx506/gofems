package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// "github.com/Josesx506/gofems/internals/api"
	"github.com/Josesx506/gofems/internals/store"
	"github.com/Josesx506/gofems/migrations"
)

type Application struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewApplication() (*Application, error) {
	// Setup DB store
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	// Migrate db from package root directory
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Stores for db access

	app := &Application{
		Logger: logger,
		DB:     pgDB,
	}

	return app, nil
}

// Add the health check controller as a method
func (a *Application) HealthChecker(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Application health is available\n")
}
