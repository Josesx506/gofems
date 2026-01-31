### Routes
Routes in go can be managed using the net/http, chi, or gorilla mux libraries. For this project, I'm using the `chi` library.


#### Setup a server
We can setup a server using the default http library and pass our chi router as one of the fields.
```go
package main

import (
    "fmt"
    "time"
    "net/http"
)

func main() {
    server := &http.Server{
        Addr:         ":3000",
        Handler:      router,
        IdleTimeout:  time.Minute,
        ReadTimeout:  10 * time.Minute,
        WriteTimeout: 30 * time.Minute,
    }

    // Update the err variable
    err = server.ListenAndServe()
    if err != nil {
        app.Logger.Fatal(err)
    }
}
```

#### Setup a basic route
Each router works like an express router and allows subroutes, nested routes and http methods. A route requires a 
- method - GET,POST,PUT ...
- route path - name of the route
- handler - similar to nodejs express controller. Includes the request object and response writer
```go
package api

import (
	"github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "git.../routes/handlers"
)

func SetupRoutes() *chi.Mux {
	router := chi.NewRouter()
    sharedHandler := &handlers.SharedHandler{}

    // Use logger middleware
	router.Use(middleware.Logger)

	router.Get("/health", sharedHandler.Health)

	return router
}
```

#### SubRoutes (Modules)
If we have a group of routes for a microservice or separate business need, we can group them together and export them as a subroute. This function takes in an existing router and 
adds routes to it. It doesn't return any variables

```go
package workouts

import (
	"github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func workoutRoutes() chi.Router {
    router := chi.NewRouter()
    workoutHandler := NewWorkoutHandler() // Same package diff file. Doesn't require import

    router.Post("/", workoutHandler.CreateWorkout)
    router.Get("/", workoutHandler.ListWorkouts)
	router.Get("/{id}", workoutHandler.GetWorkoutByID)
	router.Put("/{id}", workoutHandler.UpdateWorkoutByID)
	router.Delete("/{id}", workoutHandler.DeleteWorkoutByID)

    return router
}
```

This subroute can now be registered in the routes function as 
```go
func SetupRoutes() *chi.Mux {
	router := chi.NewRouter()
    ...
    // Ensure all workout subroutes have the prefix `workouts`
    router.Mount("/workouts", routes.workoutRoutes())
    // You can also use router.Route() to create sub routes. Check docs for details

    return router
}
```
More info can be found in the official [docs](https://go-chi.io/#/pages/routing).

#### Handlers
Handlers are the equivalent of controllers in express, flask, or fastapi. They're functions that include request and response writers enabling db calls, auth, crud operations etc. A typical handler can be defined with
```go
package api

import (
    "fmt"
    "net/http"
)

type SharedHandler struct {}

func NewSharedHandler() *SharedHandler {
    return &SharedHandler{}
}

func (sh *SharedHandler) Health(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte("Application health is available"))
    return
}
```

#### Guidelines to implement a route
1. Update database layer interface
2. Create CRUD method in database store, 
3. Create handler function
4. Add route for CRUD operation


The ***`Envelope`*** is a map with a `string` key and an `interface{}` or `any` value, allowing flexible JSON serialization by enabling different types of data to be passed, such as error messages or success data like workout structs. It provides a flexible way to structure JSON responses by allowing dynamic key-value pairs, making it easier to return different types of data with a consistent response format

`json.MarshalIndent()` differs from the standard json marshaling because it allows for indentation of JSON fields by specifying a prefix and indent, making the JSON output more readable compared to compact JSON

### Handling Errors
Log detailed error messages for internal debugging, use appropriate HTTP status codes (like StatusBadRequest or StatusInternalServerError), and provide generic error messages to clients that don't expose sensitive internal details.

- `StatusBadRequest` (400) for client-side errors like invalid input, 
- `StatusInternalServerError` (500) for server-side errors like database failures, and 
- `StatusOK` (200) for successful operations