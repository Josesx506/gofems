# gofems
Files in the `internal` directory are usually not exposed outside the module, and are 
not meant to be used by external projects. Before creating submodules, it's helpful 
to start with placing the submodules inside the internal directory till they're fully defined.

A logger provides a structured way to output messages, diagnose errors, and perform 
print debugging across the application. It can include additional information like timestamps.

- IdleTimeout (maximum time to wait for next request when Keepalives are enabled), 
- ReadTimeout (maximum duration for reading entire request), and 
- WriteTimeout (maximum time for writing response)

When creating a route, it's helpful to pass the `http.Request` as a pointer *(`r *http.Request`)* 
because the struct is large. Also, Handlers and middleware may want to modify the request 
and those changes should persist for the full life of the request.

### Parsing command line arguments
Command line arguments can be parsed using the `flag` package. This is similar to 
`argparse` in python.
```go
import (
    "fmt"
    "flag"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse() // parse args into variables so they can be used

    fmt.Printf("We are running our api on port %d\n", port)
}
```
The flag package allows us to define command line args for our functions  e.g. 
`go run main.go -- port 5132` This overwrites the original default port value of 
8080 and the text is the descripton. `os.GetEnv()` can be used to retrieve environment 
variables for security purposes.


### Folder Structure
The folder structure is a bit convoluted compared to what I'm used to. While there's good 
separation of concerns across directories, modules are packaged within structs and tacked 
onto one another requiring extensive pointer receivers to avoid bloated copies. <br>
Controllers are termed handlers and methods are inherited within structs across the 
application. Major modules in the internal folder are
```bash
└── internals
    ├── app    # big struct that inherits logs, db models, and handlers as fields
    │   └── app.go
    ├── api              # chi router is used to associate endpoints with handlers
    │   ├── routes.go
    │   └── v1           # versioned api including all routers and endpoints
    │       ├── v1Handlers.go
    │       ├── v1Router.go
    │       └── workouts # workout service as subrouter with endpoints
    │           ├── workoutHandlers.go # controllers/handlers
    │           └── workoutRouter.go
    └── store  # db queries and table schema structs
        └── workout_store.go
```

### Live server reload
To start the server, you'll need to repeatedly run the server with `go run main.go`. This doesn't track file changes and requires manually reloading the server in dev to track changes. There's a live reload tool called [air](https://github.com/air-verse/air) that supports live reloads. Install it across the environment like goose with `go install github.com/air-verse/air@latest`, then start the server with `air`. This launches the main.go file and tracks changes interactively.


### Test queries
- `curl localhost:8080/health | jq`
- `curl localhost:8080/v1/workouts/1 | jq`
- `curl -X DELETE localhost:8080/v1/workouts/1 | jq`