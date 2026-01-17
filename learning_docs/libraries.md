If you need to import a library without exclusively calling it in go
```go
import (
    _ "github.com/jackc/pgx/v4/stdlib" // Import lib without using it
)
```
The compiler won't throw an error with the underscore. This is useful if the package is 
required as a dependency in a module but not explicitly called.

### Routers (Chi)
[Chi](https://github.com/go-chi/chi) is a lightweight, idiomatic, and composable router for building Go HTTP services that implements standard library interfaces and allows for easy routing and middleware handling. A simplified example
```go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ControllerHandler(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
}

func main() {
	r := chi.NewRouter()
    // Middleware
	r.Use(middleware.Logger)

	r.Get("/", ControllerHandler)

    server := &http.Server{
		Addr:         ":3000", //port
		Handler:      r,
	}

	server.ListenAndServe()
}
```

### Databases and Migrations (pgx and goose)
[pgx](https://github.com/jackc/pgx) is a driver for accessing postgres from go. It can be installed with `go get github.com/jackc/pgx/v4/stdlib`. Go also has a standard *database/sql* library. <br>

It's a good idea to call `db.Ping()` to ensure the DSN is valid and the server is reachable 
when setting up the db connection. <br>

Add enhanced configuration to the connection pool settings with:
```go
db.SetMaxOpenConns(), db.SetMaxIdleConns(), and db.SetConnMaxIdleTime()
```


### Environment variables (godotenv)
Load dotenv files just like python. Install with `go get github.com/joho/godotenv`
```go
package main

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

func main() {
  // Load the .env file. This looks for a file named ".env" in the current directory
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  dbUrl := os.Getenv("DATABASE_URL")
}
```

To enable migrations, we need to install the `goose` package globally within the application.
