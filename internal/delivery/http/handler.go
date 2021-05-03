package http

import (
	"fmt"
	"net/http"

	"github.com/aintsashqa/go-video-service/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	VideoUuidPattern string = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[4][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"
)

type Handler struct {
	Router     chi.Router
	Repository *repository.Container
}

func NewHandler(repository *repository.Container) *Handler {
	return &Handler{
		Router:     chi.NewRouter(),
		Repository: repository,
	}
}

func (h *Handler) Handle() http.Handler {
	h.Router.Use(middleware.Logger)

	h.Router.Get("/health-check", h.HealthCheckAction)

	h.Router.Route("/api", func(r chi.Router) {

		r.Route("/video", func(r chi.Router) {

			r.Post("/", h.UploadVideoAction)
			r.Get(fmt.Sprintf("/{uuid:%s}", VideoUuidPattern), h.FindVideoAction)
			r.Get(fmt.Sprintf("/{uuid:%s}/stream", VideoUuidPattern), h.ReadBytesVideoAction)
			r.Delete(fmt.Sprintf("/{uuid:%s}", VideoUuidPattern), h.RemoveVideoAction)
		})
	})

	return h.Router
}
