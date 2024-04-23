package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"qwe/internal/client/config"
	"qwe/internal/client/ports"

	"github.com/gorilla/mux"
)

var _ Adapter = (*Server)(nil)

type Server struct {
	service ports.IService
	server  *http.Server
	cfg     *config.HTTP
	log     *slog.Logger
}

func (a *Server) Serve() error {
	router := mux.NewRouter()

	router.HandleFunc("/", a.Echo).Methods("GET")
	router.HandleFunc("/api/v1/trips", a.CreateTrip).Methods("POST").Headers("Content-Type", "application/json")

	a.server = &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", a.cfg.PORT),
	}

	return a.server.ListenAndServe()
}

func (a *Server) Shutdown(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
}

func New(service ports.IService, cfg *config.HTTP, log *slog.Logger) *Server {
	return &Server{
		service: service,
		cfg:     cfg,
		log:     log,
	}
}
