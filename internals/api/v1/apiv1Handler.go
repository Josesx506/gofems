package apiv1

import (
	"net/http"
)

type ApiV1Handler struct{}

func NewApiV1Handler() *ApiV1Handler {
	return &ApiV1Handler{}
}

func (av1 *ApiV1Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("V1 service health is active\n"))
}
