package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"qwe/internal/driver/config"
	"qwe/internal/driver/ports"

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

	router.HandleFunc("/api/v1/capturetrip", a.CaptureTrip).Methods("POST").Headers("Content-Type", "application/json")
	router.HandleFunc("/api/v1/onposition", a.OnPosition).Methods("POST").Headers("Content-Type", "application/json")
	router.HandleFunc("/api/v1/started", a.Started).Methods("POST").Headers("Content-Type", "application/json")
	router.HandleFunc("/api/v1/ended", a.Ended).Methods("POST").Headers("Content-Type", "application/json")

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
