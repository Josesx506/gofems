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
	r.Use(middleware.Logger)
	r.Get("/", ControllerHandler)

    server := &http.Server{
		Addr:         ":3000", //port
		Handler:      r,
	}
	server.ListenAndServe()
}
```