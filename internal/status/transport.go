package status

import (
	"context"
	"net"
	"net/http"

	"github.com/yushafro/03.12.2025/internal/config"
	"github.com/yushafro/03.12.2025/pkg/middleware"
)

type server struct {
	server  *http.Server
	service *service
}

func NewServer(service *service, cfg *config.Config) *server {
	srv := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		ReadTimeout:       cfg.ReadTO,
		ReadHeaderTimeout: cfg.ReadHTO,
		WriteTimeout:      cfg.WriteTO,
		IdleTimeout:       cfg.IdleTO,
	}

	return &server{
		server:  srv,
		service: service,
	}
}

func (s *server) Start() error {
	return s.server.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *server) RegisterHandlers() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", s.RegisterLinks)
	mux.HandleFunc("GET /", s.ListLinks)

	s.server.Handler = middleware.Logging(mux)
}
