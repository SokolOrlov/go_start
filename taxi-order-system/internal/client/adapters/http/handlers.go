package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"qwe/internal/client/models"
	"qwe/internal/client/ports"
)

var _ ports.IHttp = (*Server)(nil)

// Создать поездку
func (s *Server) CreateTrip(w http.ResponseWriter, r *http.Request) {
	const op = "http.createTrip"

	log := s.log.With(slog.String("op", op))

	var t models.Trip

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Error("Error handle requst", err)
		writeError(w, err)
		return
	}

	log.Info("New requst to create Trip", slog.Any("trip", t))

	s.service.CreateTrip(r.Context(), t)

	writeJSONResponse(w, http.StatusOK, nil)
}
