package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/Josesx506/gofems/internals/app"
	"github.com/Josesx506/gofems/internals/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	defer app.DB.Close() // Close the db connections at the end

	// Create a health route manually with the stdlib
	// http.HandleFunc("/health", HealthChecker)

	// Include a route handler to work with all routes.
	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port), //port
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}

	app.Logger.Printf("We are running our api on port %d\n", port)

	// Update the err variable
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
