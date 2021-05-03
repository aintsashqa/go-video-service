package http

import (
	"net/http"
)

func (_ *Handler) HealthCheckAction(w http.ResponseWriter, _ *http.Request) {
	EmptyResponse(w)
}
