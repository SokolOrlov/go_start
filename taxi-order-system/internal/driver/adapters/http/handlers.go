package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"qwe/internal/driver/models"
	"qwe/internal/driver/ports"
	status "qwe/pkg/models"
)

var _ ports.IHttp = (*Server)(nil)

// Захватить заявку
func (s *Server) CaptureTrip(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var t models.CaptureTrip

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		s.log.Error("Error capture trips", err)
		writeError(w, err)
	}

	tripstr, _ := json.Marshal(t)

	s.log.Info("http.CaptureTrip", slog.String("captureTrip", string(tripstr)))

	err = s.service.CaptureTrip(ctx, t)

	writeJSONResponse(w, http.StatusOK, err)
}

// На позиции
func (s *Server) Started(w http.ResponseWriter, r *http.Request) {
	s.updateTripStatus(r, w, status.STARTED)
}

// Старт поездки
func (s *Server) OnPosition(w http.ResponseWriter, r *http.Request) {
	s.updateTripStatus(r, w, status.ON_POSITION)
}

// Конец поездки
func (s *Server) Ended(w http.ResponseWriter, r *http.Request) {
	s.updateTripStatus(r, w, status.ENDED)
}

// Обновить статус заявки
func (s *Server) updateTripStatus(r *http.Request, w http.ResponseWriter, st status.TripStatus) {
	ctx := r.Context()

	var t models.CaptureTrip

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		s.log.Error("Error update trip status", err)
		writeError(w, err)
	}

	tripstr, _ := json.Marshal(t)

	s.log.Info("http.updateTripStatus", slog.String("captureTrip", string(tripstr)), slog.String("status", st.String()))

	err = s.service.UpdateTripStatus(ctx, t, st)

	writeJSONResponse(w, http.StatusOK, err)
}
