package apiv1

import (
	"net/http"

	"github.com/Josesx506/gofems/internal/app"
)

type ApiV1Handler struct {
	app *app.Application
}

func NewApiV1Handler(app *app.Application) *ApiV1Handler {
	return &ApiV1Handler{app: app}
}

func (av1 *ApiV1Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("V1 service health is active\n"))
}
