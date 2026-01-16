package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger log.Logger
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &Application{
		Logger: *logger,
	}

	return app, nil
}

// Add the health check controller as a method
func (a *Application) HealthChecker(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Application health is available\n")
}
