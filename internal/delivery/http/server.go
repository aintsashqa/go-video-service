package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aintsashqa/go-video-service/internal/app"
)

type Server struct {
	http *http.Server
}

func NewServer(app *app.App, handler *Handler) *Server {
	return &Server{
		http: &http.Server{
			Addr:           fmt.Sprintf("%s:%d", app.Config.HttpConfig.Host, app.Config.HttpConfig.Port),
			Handler:        handler.Handle(),
			ReadTimeout:    app.Config.HttpConfig.ReadTimeout,
			WriteTimeout:   app.Config.HttpConfig.WriteTimeout,
			MaxHeaderBytes: app.Config.HttpConfig.MaxHeaderMBytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
